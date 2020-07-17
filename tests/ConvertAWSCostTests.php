<?php declare(strict_types=1);

use App\Models\AwsCostData;
use App\Models\Transformers\Transform;
use PHPUnit\Framework\TestCase;


final class ConvertAWSCostTests extends TestCase
{
    protected $startDate = "2020-01-01";
    protected $endDate = "2020-07-01";
    protected $missingMonth = "2020-03-01";
    /**
     * make sure loads correctly
     */
    public function testObjectLoad() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample1.json"), true );
        $awsObject = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        $this->assertInstanceOf(AwsCostData::class, $awsObject);
    }


    public function testMetaData() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample1.json"), true );
        $awsObject = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        $meta = $awsObject->meta();
        $this->assertEquals("testProject1", $meta['Project']);
        $this->assertEquals("env1", $meta['Environment']);
    }

    /**
     * check it finds correct months from the original data
     * and that all months are added in so there is no gap in 
     * the time series
     */
    public function testMonthYear() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample1.json"), true );
        $awsObject = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        
        $monthYear = $awsObject->getTimePeriods();
        // this test file is missing a month
        $this->assertEquals(5, count($monthYear));

        
    }

    /**
     * check the aws service are being returned via dimensions
     */
    public function testDimensions() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample1.json"), true );
        $awsObject = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        // test file has 20 unique dimensions (aws services)
        $dimensions = $awsObject->getDimensions();
        $this->assertEquals(20, count($dimensions));        
    }

    /**
     * check the application tags are being found
     */
    public function testTags() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample1.json"), true );
        $awsObject = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        // test file has 7 tags
        $tags = $awsObject->getTags();        
        $this->assertEquals(7, count($tags));        
    }

    /**
     * Check the services data
     */
    public function testTransformToDimension() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample1.json"), true );
        $object = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        $trans = $object->transform();
        $this->assertInstanceOf(Transform::class, $trans);
        $data = $trans
            ->toDimensionByPeriod()
            ->data();
        //find secrets mananger
        $secrets = [];
        foreach($data as $row) if($row['Service'] == "SERVICE 12") $secrets[] = $row;
        // should only be one as its grouped
        $this->assertCount(1, $secrets);
        // there was no service manager in jan, so price should be 0
        $this->assertEquals(0, $secrets[0]["2020-01"]);

    }

    /**
     * checks application
     */
    public function testTransformToTag() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample1.json"), true );
        $object = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        
        $trans = $object->transform();
        $this->assertInstanceOf(Transform::class, $trans);
        $data = $trans
            ->toTagByPeriod()
            ->data();
        
        $app = [];
        foreach($data as $row) if($row['Application'] == "TAG 4") $app[] = $row;
        // should only be one as its grouped
        $this->assertCount(1, $app);
        // didnt exist in jan, so price should be 0
        $this->assertEquals(0, $app[0]["2020-01"]);
    }


    /**
     * checks services & application
     */
    public function testTransformToDimensionAndTag() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample1.json"), true );
        $object = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        
        $trans = $object->transform();
        $this->assertInstanceOf(Transform::class, $trans);
        $data = $trans
            ->toDimensionAndTagByPeriod()
            ->data();
        $taxes = [];
        $dbs = [];
        // fetch all applications
        $apps = $object->getTags();
        // tax is godo one to check, should be present for all
        foreach($data as $row) if($row['Service'] == "TAX") $taxes[] = $row;
        // dynmodb is tagged to services, good to check costs
        foreach($data as $row) if($row['Service'] == "SERVICEONCE") $dbs[] = $row;
        
        $this->assertCount( count($apps), $taxes );
        $this->assertCount( count($apps), $dbs );

        // grab one with data
        $online = [];
        foreach($dbs as $db) if($db['Application'] == "TAG 1") $online = $db;
        // rounding
        $this->assertEquals(3.23, $online['2020-01']); 
    }

    /**
     * Overall cost in Jan for sandbox product
     */
    public function testJanTotalForSimpleByProject() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample.simple.json"), true );
        $object = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        
        $data = $object
                ->getByStartTime("2020-01-01")
                ;
        $cost = 0.0;
        foreach($data as $item) $cost += $object->getCost($item);    
        $this->assertEquals(55.63, $cost);

        $data = $object
                    ->transform()
                    ->toProjectByPeriod()
                    ->data()
                    ;
        $this->assertEquals(55.63, $data[0]['2020-01']);

    }

    /**
     * 
     */
    public function testCostByProject() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample.service-one.json"), true );
        $object = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        
        $data = $object
            ->transform()
            ->toDimensionByPeriod()
            ->data();
        
        // grab dynamodb
        $dbs = [];
        foreach($data as $row){
            if($row['Service'] == "SERVICEONCE") $dbs[] = $row;
        }
        $totaled = array_shift($dbs);
        $this->assertEquals(4.83, $totaled['2020-01']);
        $this->assertEquals(7.14, $totaled['2020-05']);
        
    }

    /**
     * 
     */
    public function testJanEnvironmentCostForSimpleProject() : void
    {
        $testData = json_decode( file_get_contents(__DIR__ ."/data/sample.simple.json"), true );
        $object = new AwsCostData(
            $testData['data'], 
            $testData['startDate'], 
            $testData['endDate'], 
            $testData['environment'], 
            $testData['project']
        );
        
        $data = $object
            ->transform()
            ->toEnvironmentByPeriod()
            ->data();
        // is the same as by project - as just a single env and project
        $this->assertEquals(55.63, $data[0]['2020-01']);
        $this->assertEquals(63.85, $data[0]['2020-02']);
        
    }


}
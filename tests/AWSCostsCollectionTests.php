<?php declare(strict_types=1);

use App\Models\AwsCostData;
use App\Models\AwsCostDataCollection;
use App\Models\Transformers\Transform;
use PHPUnit\Framework\TestCase;


final class AWSCostsCollectionTests extends TestCase
{
    
    /**
     * make sure loads correctly
     */
    public function testObjectLoad() : void
    {
        // single data file
        $collection = new AwsCostDataCollection([
            __DIR__ ."/data/sample.service-one.json"
        ]);
        $this->assertInstanceOf(AwsCostDataCollection::class, $collection);
    }

    /**
     * Check the loading of environments
     */
    public function testEnvironments() : void
    {
        // data files
        $sample1 = __DIR__ ."/data/sample1.json";
        $sample2 = __DIR__ ."/data/sample2.json";
        $sample3 = __DIR__ ."/data/sample3.json";
        // single data file
        $collection = new AwsCostDataCollection([
            $sample1
        ]);
        
        $byEnv = $collection->getPerMonth('Environment');
        $this->assertCount(1, $byEnv);
        // add the same file again
        $collection->addCost( $sample1 );
        // should still only be a single env & no data change
        $byEnv = $collection->getPerMonth('Environment');
        $this->assertCount(1, $byEnv);
        // jan cost should be...
        $janCost = 55.63;
        $this->assertEquals($janCost, $byEnv[0]['2020-01']);
        // add another project, same costs with different name
        $collection->addCost( $sample3 );
        $byEnv = $collection->getPerMonth('Environment');
        $this->assertEquals( $janCost * 2, $byEnv[0]['2020-01']);
        // add another cost file, but has new env name
        $collection->addCost( $sample2 );
        $byEnv = $collection->getPerMonth('Environment');
        $this->assertCount(2, $byEnv);
        // env1
        $env1 = false;
        foreach($byEnv as $env) if($env['Environment'] == "env1") $env1 = $env;
        $this->assertEquals( $janCost * 2, $env1['2020-01']);

        $env2 = false;
        foreach($byEnv as $env) if($env['Environment'] == "env2") $env2 = $env;
        $this->assertEquals( $janCost, $env2['2020-01']);

    }

    /**
     * 
     */
    public function testProjects() : void
    {
        $janCost = 55.63;
        // data files
        $sample1 = __DIR__ ."/data/sample1.json";
        $sample2 = __DIR__ ."/data/sample2.json";
        $sample3 = __DIR__ ."/data/sample3.json";
        // single data file
        $collection = new AwsCostDataCollection([
            $sample1,
            $sample2, 
            $sample3
        ]);
        
        $byProj = $collection->getPerMonth('Project');
        // sample1 & sample2 have the same project name
        $this->assertCount(2, $byProj);
        
        $proj1 = false;
        foreach($byProj as $p) if($p['Project'] == "testProject1") $proj1 = $p;
        $this->assertEquals( $janCost * 2, $proj1['2020-01']);

        $proj2 = false;
        foreach($byProj as $p) if($p['Project'] == "testProject2") $proj2 = $p;
        $this->assertEquals( $janCost, $proj2['2020-01']);

    }

    /**
     * 
     */
    public function testServices() : void
    {
        // costs are from the aws console for sandbox
        $junCost = 9.98;
        // data files
        $sample  = __DIR__ ."/data/sample.service-one.json";
        $sample1 = __DIR__ ."/data/sample1.json";
        $sample2 = __DIR__ ."/data/sample2.json";
        $sample3 = __DIR__ ."/data/sample3.json";
        // single data file
        $collection = new AwsCostDataCollection([
            $sample
        ]);        
        $by = $collection->getPerMonth('Service');
        $this->assertCount(1, $by);
        $dbs = [];
        foreach($by as $row){
            if($row['Service'] == "SERVICEONCE") $dbs[] = $row;
        }
        $totaled = array_shift($dbs);
        $this->assertEquals(4.83, $totaled['2020-01']);
        $this->assertEquals(7.14, $totaled['2020-05']);


        $collection = new AwsCostDataCollection([
            $sample1,
            $sample2,
            $sample3,
        ]);        
        $by = $collection->getPerMonth('Service');
        // 20 services
        $this->assertCount(20, $by);
        $found = false;
        
        foreach($by as $i) if($i['Service'] == "TAX") $found = $i;
        $this->assertEquals( $junCost * 3, $found['2020-06']);

    
    }

    /**
     * 
     */
    public function testApplications() : void
    {
        // costs are from the aws console for sandbox
        $expectedCost = 3.03;
        $depsCost = 0.11;
        // data files
        $singleservice  = __DIR__ ."/data/sample.service-one.json";
        $sample1 = __DIR__ ."/data/sample1.json";
        $sample2 = __DIR__ ."/data/sample2.json";
        $sample3 = __DIR__ ."/data/sample3.json";
        // single data file
        $collection = new AwsCostDataCollection([
            $singleservice
        ]);        
        $by = $collection->getPerMonth('Application');
        // only two in the dynamodb file
        $this->assertCount(2, $by);
        $found = false;

        foreach($by as $i) if($i['Application'] == "TAG 1") $found = $i;
        $this->assertEquals( $expectedCost, $found['2020-02']);
        
        $collection = new AwsCostDataCollection([
            $sample1
        ]);        

        $by = $collection->getPerMonth('Application');
        
        // check - only has a cost in march
        foreach($by as $i) if($i['Application'] == "TAG 2") $found = $i;
        
        $this->assertEquals( 0, $found['2020-01']);
        $this->assertEquals( $depsCost, $found['2020-02']);
        $this->assertEquals( 0, $found['2020-03']);
        $this->assertEquals( 0, $found['2020-04']);
        $this->assertEquals( 0, $found['2020-05']);
        $this->assertEquals( 0, $found['2020-06']);

        $collection = new AwsCostDataCollection([
            $sample1,
            $sample2,
            $sample3
        ]); 
        $by = $collection->getPerMonth('Application');
        
        // check  - only has a cost in march
        foreach($by as $i) if($i['Application'] == "TAG 2") $found = $i;
        
        $this->assertEquals( $depsCost * 3, $found['2020-02']);
    }





    public function testProjectService() : void
    {
        // costs are from the aws console for sandbox
        $expectedCost = 4.83;
        // data files
        $sample  = __DIR__ ."/data/sample.service-one.json";
        // single data file
        $collection = new AwsCostDataCollection([
            $sample
        ]);        
        $by = $collection->getPerMonth(['Project', 'Service']);
        //  has one project, one service, should match the total
        $found = false;
        foreach($by as $i) if($i['Service'] == "SERVICEONCE") $found = $i;
        $this->assertEquals( $expectedCost, $found['2020-01']);

        
    }


}
<?php declare(strict_types=1);

use Aws\CostExplorer\CostExplorerClient;
use Aws\Sts\StsClient;
use PHPUnit\Framework\TestCase;


final class AWSIdentityTests extends TestCase
{
    
    public function testIdentity() : void
    {
        $arns = [
            //sandbox
            [
                'arn' => "arn:aws:iam::995199299616:role/breakglass",
                'month' => '2020-02-01',
                'value' => '63.85'
            ],
            //sirius preprod
            [
                'arn' => "arn:aws:iam::492687888235:role/breakglass",
                'month' => '2020-03-01',
                'value' => '8897.45'
            ]
        ];

        $stsClient = new StsClient([        
            'region' => 'eu-west-1',
            'version' => 'latest'
        ]);
        
        foreach($arns as $test)
        {
            // assume the role
            $role = $stsClient->AssumeRole([
                'RoleArn' => $test['arn'],
                'RoleSessionName' => "get-costs-cli",
                'DurationSeconds' => 900
            ]);
            // create the client
            $client = new CostExplorerClient([
                'region'        => 'eu-west-1', 
                'version'       => 'latest',
                'credentials' =>  [
                    'key'    => $role['Credentials']['AccessKeyId'],
                    'secret' => $role['Credentials']['SecretAccessKey'],
                    'token'  => $role['Credentials']['SessionToken']
                ]
            ]);
            // get results
            $result = $client->getCostAndUsage([ 
                    'Metrics' => [ 'BlendedCost' ],
                    'Granularity' => 'MONTHLY',
                    'TimePeriod' => [ 'Start' => "2020-01-01", 'End' => "2020-07-01"]
            ]);
            // get the data
            $data = $result->get("ResultsByTime");
            // find the month value from the data
            $foundValue = 0.0;
            foreach($data as $month)
            {
                if($month['TimePeriod']['Start'] == $test['month'])
                {
                    $foundValue = round( $month['Total']['BlendedCost']['Amount'], 2);
                    $this->assertEquals($test['value'], $foundValue);
                }
            }
            // make sure its found a value
            $this->assertGreaterThan(1, $foundValue);

        }

    }


}

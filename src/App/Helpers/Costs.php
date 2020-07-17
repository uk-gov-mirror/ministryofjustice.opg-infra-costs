<?php
namespace App\Helpers;

use Aws\CostExplorer;
use Aws\CostExplorer\CostExplorerClient;
use Aws\Credentials\CredentialProvider;

class Costs
{
    
    /**
     * Fetches cost data grouped by the service type and application tag
     * 
     */
    public static function blendedGroupedByServiceAndTag(
        string $start, 
        string $end, 
        string $granularity = 'MONTHLY'
        ) 
    {
        $client = new CostExplorerClient([
            'region'        => 'eu-west-1', 
            'version'       => '2017-10-25',
            'credentials'   => CredentialProvider::env()
        ]);
        $result = $client->getCostAndUsage(
            [
                'Metrics' => [
                    'BlendedCost'
                ],
                'Granularity' => $granularity,
                'TimePeriod' => [
                    'Start' => $start,
                    'End' => $end
                ],
                'GroupBy' => [
                    [
                        'Key' => 'SERVICE',
                        'Type' => 'DIMENSION'
                    ],
                    [
                        'Key' => 'application',
                        'Type' => 'TAG'
                    ]
                ]
            ]
        );
        $data = $result->get("ResultsByTime");
        return $data;
    }



    
    

    
}
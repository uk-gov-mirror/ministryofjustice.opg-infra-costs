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
        CostExplorerClient $client,
        string $granularity = 'MONTHLY'
        ) 
    {
        
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

<?php
namespace App\Helpers;

use Aws\CostExplorer\CostExplorerClient;
use Aws\Sts\StsClient;

class AssumedRoleClient
{
    
    /**
     * Fetches a client using mfa for the arn
     *
     */
    public static function get(       
        string $arn
        ) : CostExplorerClient
    {
        $stsClient = new StsClient([        
            'region' => 'eu-west-1',
            'version' => 'latest'
        ]);

         // assume the role
         $role = $stsClient->AssumeRole([
            'RoleArn' => $arn,
            'RoleSessionName' => "get-costs-cli",
            'DurationSeconds' => 900
        ]);
        // create the client
        return new CostExplorerClient([
            'region'        => 'eu-west-1', 
            'version'       => 'latest',
            'credentials' =>  [
                'key'    => $role['Credentials']['AccessKeyId'],
                'secret' => $role['Credentials']['SecretAccessKey'],
                'token'  => $role['Credentials']['SessionToken']
            ]
        ]);
        
        
    }
    
}

<?php
namespace App\Helpers;

use Aws\CostExplorer\CostExplorerClient;
use Aws\Sts\StsClient;

class AssumedRoleClient
{
    
    /**
     * Fetches a client using mfa for the arn
     */
    public static function get(
        string $identityAccount,
        string $identityUser,
        string $mfaToken,
        string $arn
        ) : CostExplorerClient
    {
        $stsClient = new StsClient([        
            'region' => 'eu-west-1',
            'version' => 'latest'
        ]);

        $result = $stsClient->AssumeRole([
            'RoleArn' => $arn,
            'RoleSessionName' => "get-costs-cli",
            'SerialNumber' => "arn:aws:iam::${identityAccount}:mfa/${identityUser}",
            'TokenCode' => $mfaToken
        ]);


        return new CostExplorerClient([
            'region'        => 'eu-west-1', 
            'version'       => 'latest',
            'credentials' =>  [
                'key'    => $result['Credentials']['AccessKeyId'],
                'secret' => $result['Credentials']['SecretAccessKey'],
                'token'  => $result['Credentials']['SessionToken']
            ]
        ]);
        
    }



    
    

    
}

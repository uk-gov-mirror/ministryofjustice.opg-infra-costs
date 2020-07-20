<?php
namespace App\Helpers;

use Aws\CostExplorer\CostExplorerClient;
use Aws\Sts\StsClient;

class AssumedRoleClient
{
    protected static $result = null;
    
    /**
     * Fetches a client using mfa for the arn
     *
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

        if(self::$result == null)
        {
            self::$result = $stsClient->AssumeRole([
                'RoleArn' => $arn,
                'RoleSessionName' => "get-costs-cli",
                'SerialNumber' => "arn:aws:iam::${identityAccount}:mfa/${identityUser}",
                'TokenCode' => $mfaToken
            ]);
        }
        
        return new CostExplorerClient([
            'region'        => 'eu-west-1', 
            'version'       => 'latest',
            'credentials' =>  [
                'key'    => self::$result['Credentials']['AccessKeyId'],
                'secret' => self::$result['Credentials']['SecretAccessKey'],
                'token'  => self::$result['Credentials']['SessionToken']
            ]
        ]);
        
    }



    
    

    
}

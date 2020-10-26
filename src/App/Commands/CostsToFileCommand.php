<?php
namespace App\Commands;

use App\Helpers\AssumedRoleClient;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

use App\Helpers\Costs;

class CostsToFileCommand extends Command
{
    protected string $mfaToken = "";
    protected static string $identityAccountId = "631181914621";
    protected $arns = [
        'sandbox-sandbox-breakglass'        => "arn:aws:iam::995199299616:role/breakglass",
        'identity-identity-breakglass'      => "arn:aws:iam::631181914621:role/breakglass",
        // sirius
        'sirius-dev-breakglass'             => "arn:aws:iam::288342028542:role/breakglass",
        'sirius-preprod-breakglass'         => "arn:aws:iam::492687888235:role/breakglass",
        'sirius-prod-breakglass'            => "arn:aws:iam::649098267436:role/breakglass",
        'sirius-backup-breakglass'          => "arn:aws:iam::132068124730:role/breakglass",
        // serve
        'serve-dev-breakglass'              => "arn:aws:iam::705467933182:role/breakglass",
        'serve-preprod-breakglass'          => "arn:aws:iam::540070264006:role/breakglass",
        'serve-prod-breakglass'             => "arn:aws:iam::933639921819:role/breakglass",
        // lpa
        'lpa-dev-breakglass'                => "arn:aws:iam::050256574573:role/breakglass",
        'lpa-preprod-breakglass'            => "arn:aws:iam::987830934591:role/breakglass",
        'lpa-prod-breakglass'               => "arn:aws:iam::980242665824:role/breakglass",
        // digideps
        'digideps-dev-breakglass'           => "arn:aws:iam::248804316466:role/breakglass",
        'digideps-preprod-breakglass'       => "arn:aws:iam::454262938596:role/breakglass",
        'digideps-prod-breakglass'          => "arn:aws:iam::515688267891:role/breakglass",
        // refunds
        'refunds-dev-breakglass'            => "arn:aws:iam::936779158973:role/breakglass",
        'refunds-preprod-breakglass'        => "arn:aws:iam::764856231715:role/breakglass",
        'refunds-prod-breakglass'           => "arn:aws:iam::805626386523:role/breakglass",
        // ual
        'ual-dev-breakglass'                => "arn:aws:iam::367815980639:role/breakglass",
        'ual-preprod-breakglass'            => "arn:aws:iam::888228022356:role/breakglass",
        'ual-prod-breakglass'               => "arn:aws:iam::690083044361:role/breakglass",
        // org
        'org-management-breakglass'         => "arn:aws:iam::311462405659:role/breakglass",
        // old roles, but actively used 
        'jenkins-dev-accountwrite'          => "arn:aws:iam::679638075911:role/account-write",
        'jenkins-prod-accountwrite'         => "arn:aws:iam::997462338508:role/account-write",
        // LEGACY
        'refunds-LEGACY_dev-accountwrite'   => "arn:aws:iam::792093328875:role/account-write",
        'refunds-LEGACY_prod-accountwrite'  => "arn:aws:iam::574983609246:role/account-write",
        'lpa-LEGACY_prod-breakglass'        => "arn:aws:iam::550790013665:role/breakglass",
        'org-shared-accountwrite'           => "arn:aws:iam::357766484745:role/account-write",
    ];

    protected function configure()
    {
        $this->setName("costs-to-file")
            ->setDescription("Download costs to a local file")            
            ->addOption(
                "startDate", 
                "s", 
                InputOption::VALUE_OPTIONAL, 
                "Start Date for the query - YYYY-MM-DD",
                date("Y-m-01", mktime(0,0,0, 1, 1, date("Y") ) )
                )
            ->addOption(
                "endDate", 
                "e", 
                InputOption::VALUE_OPTIONAL, 
                "End Date for the query - YYYY-MM-DD",
                // first day of next month
                date("Y-m-01", mktime(0,0,0, date("n")+1, 1, date("Y") ) )
                )
            ->addOption(
                "awsIdentityAccountId",
                'i',
                InputOption::VALUE_OPTIONAL,
                "The account id to use for identity MFA arn",
                self::$identityAccountId
            )
            ;
    }

    

    

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln("<info>Fetching costs...</info>");
        $start = $input->getOption("startDate");
        $end = $input->getOption("endDate");
        
        $identityAccount = $input->getOption("awsIdentityAccountId");
        
        $counter = 1;
        $max = count($this->arns);
        foreach($this->arns as $name => $arn){
            list($project, $environment, $role) = explode("-", $name);
            $output->writeln("<info>Getting data [${counter}/${max}] for ${name}</info>");
            
            // assumed role
            $client = AssumedRoleClient::get($arn);

            $data = Costs::blendedGroupedByServiceAndTag($start, $end, $client);
            $page = [
                'environment' => $environment,
                'project' => $project,
                'startDate' => $start,
                'endDate' => $end,
                'data' => $data
            ];
            $file = __DIR__ . "/../../../tmp/".$project. "." . $environment. ".json";
            $dir = dirname($file);
            if(!is_dir($dir)) mkdir($dir, 0777, true);
            file_put_contents($file, json_encode($page));
            $counter ++ ;
        }
        return Command::SUCCESS;
        
        
    }
    
}

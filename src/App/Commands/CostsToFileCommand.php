<?php
namespace App\Commands;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

use App\Helpers\Identity;
use App\Helpers\Costs;

class CostsToFileCommand extends Command
{
    protected $project = "";
    protected $environment = "";
    protected $role = "";

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
            ;
    }

    

    protected function awsRole(OutputInterface $output)
    {
        // fetch role name from env
        $vaultRole = getenv("AWS_VAULT");
        if(! $vaultRole)
        {
            $output->writeln("<error>Cannot find AWS Vault role</error>"); 
            return false;
        }

        list($this->project, $this->environment, $this->role) = explode("-", $vaultRole);        
        $output->writeln("<info>Using vault role: ${vaultRole} =></info>");
        $output->writeln("<info>- project: ". $this->project ."</info>");
        $output->writeln("<info>- environment: ". $this->environment ."</info>");
        $output->writeln("<info>- role: ".$this->role."</info>");
        return true;
    }


    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln("<info>Fetching costs...</info>");

        if(! $this->awsRole($output) )
        {
            return Command::FAILURE;
        }        

        if($this->project && $this->environment && $this->role)
        {
            $start = $input->getOption("startDate");
            $end = $input->getOption("endDate");
            $data = Costs::blendedGroupedByServiceAndTag($start, $end);
            $page = [
                'environment' => $this->environment,
                'project' => $this->project,
                'startDate' => $start,
                'endDate' => $end,
                'data' => $data
            ];
            $dir = __DIR__ . "/../../../tmp/".$this->project. "." . $this->environment. ".json";
            file_put_contents($dir, json_encode($page));
            return Command::SUCCESS;
        }
        else
        {
            $output->writeln("<error>Project / Environment / Role not found from ENV VARS</error>");
            return Command::FAILURE;
        }
        
    }
    
}
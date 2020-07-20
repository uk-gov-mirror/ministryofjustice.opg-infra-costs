<?php
namespace App\Commands;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;
use PhpOffice\PhpSpreadsheet\Writer\Xlsx;

use App\Helpers\CostSpreadsheet;
use App\Helpers\DrawChart;
use App\Models\AwsCostDataCollection;
use PhpOffice\PhpSpreadsheet\Chart\DataSeries;

class GenerateSpreadsheetCommand extends Command
{
    
    protected function configure()
    {
        $this->setName("spreadsheet")
            ->setDescription("Convert series of json export files to a singular spreadsheet for viewing")            
            ->addOption(
                "directory", 
                "d", 
                InputOption::VALUE_OPTIONAL, 
                "Directory of data files to load",
                "./tmp/"                
            )
            ->addOption(
                "file", 
                "f", 
                InputOption::VALUE_OPTIONAL, 
                "Location to save spreadsheet to",
                "./costs.xlsx"                
            )           
            ;
    }

    

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        
        $output->writeln("<info>Creating spreadsheet...</info>");

        $dir = $input->getOption("directory");
        $files = glob($dir."*.json");
        $collection = new AwsCostDataCollection($files);

        $dataset =[
            'Project' => [
                'data' => $collection->getPerMonth('Project'),
                'chart' => DataSeries::TYPE_BARCHART
            ],
            'Service' => [
                'data' => $collection->getPerMonth('Service'),
                'chart' => DataSeries::TYPE_BARCHART
            ],
            'Environment' => [
                'data' => $collection->getPerMonth('Environment'),
                'chart' => DataSeries::TYPE_BARCHART
            ] ,
            // not widely used enough at the moment
            'Application' => [
                'data' => $collection->getPerMonth('Application')
            ],
            'ProjectPerService' => [
                'data' => $collection->getPerMonth(['Project', 'Service'])
            ],
            'ProjectPerEnvironment' => [
                'data' => $collection->getPerMonth(['Project', 'Environment'])                
            ],
            // too many data points for excel
            'ProjectPerEnvironmentPerService' => [
                'data'=> $collection->getPerMonth(['Project', 'Environment','Service'])
            ],

            'Totals' => [
                'data' => $collection->getTotalPerMonth(),
                'chart' => DataSeries::TYPE_BARCHART
            ],
            
        ];
        
        $file = $input->getOption("file");
        $spreadsheet = CostSpreadsheet::getSpreadsheet($file);

        // do this first as it creates a file lock / race
        foreach($dataset as $worksheetName => $sheet) CostSpreadsheet::writeToWorksheet($sheet['data'], $worksheetName, $spreadsheet, false);  
        
        foreach($dataset as $worksheetName => $sheet)
        {            
            if(isset($sheet['chart']))
            {
                DrawChart::draw($sheet['chart'], $sheet['data'], $worksheetName, $worksheetName, $spreadsheet);
                
            }
        }
        $spreadsheet->setActiveSheetIndex(0);
        $writer = new Xlsx($spreadsheet);
        $writer->setIncludeCharts(true);
        $writer->save($file);

        return Command::SUCCESS;
        
        
    }
    
}

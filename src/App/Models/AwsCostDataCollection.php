<?php
namespace App\Models;
use App\Models\Transformers\Transform;

class AwsCostDataCollection
{
    protected $awsCosts;
    protected $sourceFiles;
    protected $endDate;

    public function __construct($sourceFilesArray)
    {
        $this->sourceFiles = $sourceFilesArray;
        $this->load();
    }

    /**
     * Add new costs to the set
     */
    public function addCost($filePath) : void
    {
        $this->sourceFiles[] = $filePath;
        $this->sourceFiles = array_unique($this->sourceFiles);
        $this->load();
    }
    /**
     * generate the array of cost entities
     */
    protected function load()
    {
        $this->awsCosts = [];
        foreach($this->sourceFiles as $filePath)
        {
            $data = json_decode( file_get_contents($filePath), true);    
            $this->awsCosts[] = new AwsCostData(
                $data['data'], 
                $data['startDate'], 
                $data['endDate'], 
                $data['environment'], 
                $data['project']
            );
        }
    }
    
    /**
     * Overarching totals
     */
    public function getTotalPerMonth()
    {
        $flatData = [];
        foreach($this->awsCosts as $cost)
        {
            $flat = $cost->transform()->toByPeriod()->data();
            // merge in
            foreach($flat as $date => $value)
            {
                if(!isset($flatData[$date])) $flatData[$date] = 0.0;
                $flatData[$date] += $value;
            }
            
        }
        return [$flatData];
    }


    public function getPerMonth($what)
    {
        
        // if array, sort it join elements
        if(is_array($what))
        {
            $asked = $what;
            sort($what);
            $what = implode($what);

        } else $asked = [$what];

        $function = null;
        switch($what) {
            /**
             * Multi key
             */
            case "ProjectService":
                $function = "toProjectAndDimensionByPeriod";
                break;
            case "EnvironmentProject":
                $function = "toProjectAndEnvironmentByPeriod";
                break;
            case "EnvironmentProjectService":
                $function = "toProjectAndEnvironmentAndDimensionByPeriod";
                break;
            /**
             * Single key versions
             */
            case "Environment": 
                $function = "toEnvironmentByPeriod";
                break;
            case "Project": 
                $function = "toProjectByPeriod"; 
                break;
            case "Service": 
                $function = "toDimensionByPeriod"; 
                break;
            case "Application":
                $function = "toTagByPeriod"; 
                break;
            default:
                throw new \Exception("Cannont get per month : $what : unknown type");
                break;
        }
        $flatData = $this->transformed($function);
        return $this->summation($flatData, $asked);

    }

    /**
     * Handle the call for each awsCost item to the transformer
     * and fetch the data from that, merge it together in a 
     * single flat array
     */
    public function transformed(string $functionName)
    {
        $flatData = [];
        foreach($this->awsCosts as $cost)
        {
            $flatData = array_merge(
                        $flatData, 
                        $cost->transform()->$functionName()->data()
                );
        }
        return $flatData;
    }

    /**
     * using the flat array, find matching records based on the value
     * of $col and 
     */
    public function summation($flatData, $cols)
    {
        $data = [];
        $result = [];
        foreach($flatData as $row)
        {
            $key = $this->key($row, $cols);
            if(! isset($data[$key]) ) $data[$key] = $row;
            else{
                foreach($row as $k => $v)
                {
                    //only add numerics
                    if(is_numeric($v) ) $data[$key][$k] += $v;
                }
            }
        }
        foreach($data as $row) $result[] = $row;
        return $result;
    }


    protected function key($row, $cols) : string
    {
        $key = "";
        foreach($cols as $col) $key .= $row[$col];
        return $key;
    }
}
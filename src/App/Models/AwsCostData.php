<?php
namespace App\Models;
use App\Models\Transformers\Transform;

class AwsCostData
{
    protected $startDate;
    protected $endDate;

    protected $environment = "dev";
    protected $project = "proj";

    public function __construct(
        $costData, 
        string $startDate, 
        string $endDate, 
        string $env = "dev", 
        string $proj = "pro")
    {
        $this->startDate = $startDate;
        $this->endDate = $endDate;
        $this->costs = $costData;
        $this->environment = $env;
        $this->project = $proj;
    }

    /**
     * return meta info
     */
    public function meta()
    {
        return [
            'Project' => $this->project,
            'Environment' => $this->environment,
            'startDate' => $this->startDate,
            'endDate' => $this->endDate
        ];
    }

    /**
     * return all the raw data
     */
    public function data()
    {
        return $this->costs;
    }
    /**
     * Find all the timePeriods, just use the start date as we
     * only do monthly
     */
    public function getTimePeriods()
    {
        $times = [];
        foreach($this->costs as $data)
        {
            $times[] = date("Y-m", strtotime($data['TimePeriod']['Start']) );
        }
        array_unique($times);
        return $times;
    }

    /**
     * Find all possible values for
     * - DIMENSION     
     */
    public function getDimensions()
    {
        $set = [];
        foreach($this->costs as $data)
        {
            $groups = $data['Groups'];
            foreach($groups as $group)
            {
                $set[] = strtoupper(array_shift($group["Keys"]) );
            }
        }
        $set = array_unique($set);
        sort($set);
        return $set;
    }
    /**
     * Find all possible values for
     * - TAG     
     */
    public function getTags()
    {
        $set = [];
        foreach($this->costs as $data)
        {
            $groups = $data['Groups'];
            foreach($groups as $group)
            {
                $tag = str_ireplace("APPLICATION$", "", strtoupper(array_pop($group["Keys"]) ) );
                if($tag == "") $tag = "NA";
                $set[] = $tag;
            }
        }
        $set = array_unique($set);
        sort($set);
        return $set;
    }

    /**
     * 
     */
    public function getByStartTime($period)
    { 
        $found = [];
        foreach($this->costs as $data)
        {
            $time = $data['TimePeriod']['Start'];
            if($time === $period) $found[] = $data;
        }
        return $found;

    }
    /**
     * 
     */
    public function getCostByDimensionAndTag($set, $dimension, $tag) : float
    {
        if(!isset($set['Groups'])) return false;  
        $cost = 0.0;      
        foreach($set['Groups'] as $group)
        {
            $t = str_ireplace("APPLICATION$", "", strtoupper($group["Keys"][1] ) );
            if($t == "") $t = "NA";

            if(strtoupper($group['Keys'][0]) == $dimension && $t == $tag){
                $cost += floatval($group['Metrics']['BlendedCost']['Amount']);
            }
        }
        return round($cost, 2);
    }
    /**
     * 
     */
    public function getCostByDimension($set, $dimension) : float
    {
        if(!isset($set['Groups'])) return false;      
        $cost = 0.0;
        foreach($set['Groups'] as $group)
        {
            if(strtoupper($group['Keys'][0]) == $dimension){
                $cost += floatval($group['Metrics']['BlendedCost']['Amount']);
            }
        }
        return round($cost, 2);
    }

    public function getCostByTag($set, $tag) : float
    {
        if(!isset($set['Groups'])) return false;    
        $cost = 0.0;    
        foreach($set['Groups'] as $group)
        {
            $t = str_ireplace("APPLICATION$", "", strtoupper($group["Keys"][1] ) );
            if($t == "") $t = "NA";
            if($t == $tag){
                $cost += floatval($group['Metrics']['BlendedCost']['Amount']);
            }
        }
        return round($cost, 2);
    }

    /**
     * 
     */
    public function getCost($set) : float
    {
        if(!isset($set['Groups'])) return false;
        $cost = 0.0;
        foreach($set['Groups'] as $group)
        {
            $val = floatval($group['Metrics']['BlendedCost']['Amount']);
            $cost += $val;
            
        }
        return round( $cost, 2);
    }

    /**
     * Find months that are missing from $known
     */
    protected function missingMonths($start, $end, $known)
    {
        $range = [];
        $current = strtotime($start);
        $endTime = strtotime($end);
        while($current < $endTime)
        {
            $range[] = date("Y-m", $current);
            $current = strtotime("+1 Month", $current);
        }
        return array_diff($range, $known);        
    }


    /**
     * Return the transformation class
     */
    public function transform()
    {
        return new Transform($this->costs, $this);
    }

    /**
     * Find missing months in the TimePeriod
     * and fill in the missing parts
     */
    public function dateRange()
    {
        $dates = $this->getTimePeriods();
        $missing = $this->missingMonths(
                            $this->startDate, 
                            $this->endDate, 
                            $dates
        );
        $dates = array_merge($dates, $missing);        
        sort($dates);        
        return $dates;
    }


}
<?php
namespace App\Models\Transformers;

class ByPeriod
{
    protected $data;
    public function __construct($data, $aws)
    {
        $this->data = $data;
        $this->aws = $aws;
    }

    public function data()
    {
        $structure = [];
        $meta = $this->aws->meta();
        $periods = $this->aws->dateRange();
        // convert and map to array strucure of service per row
        foreach($periods as $date)
        {
            $structure[$date] = 0.0;
            $forDate = $this->aws->getByStartTime($date . "-01");
            foreach($forDate as $item)
            {
                $cost = $this->aws->getCost($item);                
                $structure[$date] = $structure[$date] + $cost;
            }
        }

        return $structure;
    }

    

}
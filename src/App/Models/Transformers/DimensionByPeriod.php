<?php
namespace App\Models\Transformers;

class DimensionByPeriod
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
        $dimensions = $this->aws->getDimensions();
        $periods = $this->aws->dateRange();
        // convert and map to array strucure of service per row
        foreach($dimensions as $service)
        {
            $struct = [                
                'Service' => $service
            ];
            foreach($periods as $date)
            {
                $struct[$date] = 0.0;
                $forDate = $this->aws->getByStartTime($date . "-01");

                foreach($forDate as $item)
                {
                    $cost = $this->aws->getCostByDimension($item, $service);                    
                    $struct[$date] = $struct[$date] + $cost;
                }
            }
            $structure[] = $struct;
        }

        return $structure;
    }

    

}
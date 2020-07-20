<?php
namespace App\Models\Transformers;

class TagByPeriod
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
        $tags = $this->aws->getTags();
        $periods = $this->aws->dateRange();
        // convert and map to array strucure of service per row
        foreach($tags as $tag)
        {
            $struct = [
                'Application' => $tag
            ];
            foreach($periods as $date)
            {
                $struct[$date] = 0.0;
                $forDate = $this->aws->getByStartTime($date . "-01");

                foreach($forDate as $item)
                {
                    $cost = $this->aws->getCostByTag($item, $tag);                    
                    $struct[$date] = $struct[$date] + $cost;
                }
            }
            $structure[] = $struct;
        }

        return $structure;
    }

    

}
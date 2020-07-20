<?php
namespace App\Models\Transformers;

class ProjectAndDimensionByPeriod extends DimensionByPeriod
{
   

    public function data()
    {
        $meta = $this->aws->meta();
        $data = parent::data();
        $structure = [];
        foreach($data as $d)
        {
            $structure[] = array_merge(['Project' => $meta['Project']], $d);
        }

        return $structure;
    }


    

    

}
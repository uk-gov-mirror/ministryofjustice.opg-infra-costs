<?php
namespace App\Models\Transformers;

class Transform
{
    protected $data;
    protected $aws;
    public function __construct($data, $aws)
    {
        $this->data = $data;
        $this->aws = $aws;
    }

    public function toDimensionByPeriod() : DimensionByPeriod
    {
        return new DimensionByPeriod($this->data, $this->aws);
    }

    public function toTagByPeriod() : TagByPeriod
    {
        return new TagByPeriod($this->data, $this->aws);
    }

    public function toDimensionAndTagByPeriod() : DimensionAndTagByPeriod
    {
        return new DimensionAndTagByPeriod($this->data, $this->aws);
    }

    public function toProjectByPeriod() : ProjectByPeriod
    {
        return new ProjectByPeriod($this->data, $this->aws);
    }
    
    public function toEnvironmentByPeriod() : EnvironmentByPeriod
    {
        return new EnvironmentByPeriod($this->data, $this->aws);
    }

    public function toProjectAndDimensionByPeriod() : ProjectAndDimensionByPeriod
    {
        return new ProjectAndDimensionByPeriod($this->data, $this->aws);
    }
    
    public function toProjectAndEnvironmentByPeriod() : ProjectAndEnvironmentByPeriod
    {
        return new ProjectAndEnvironmentByPeriod($this->data, $this->aws);
    }

    public function toProjectAndEnvironmentAndDimensionByPeriod() : ProjectAndEnvironmentAndDimensionByPeriod
    {
        return new ProjectAndEnvironmentAndDimensionByPeriod($this->data, $this->aws);
    }

    public function toByPeriod() : ByPeriod
    {
        return new ByPeriod($this->data, $this->aws);
    }
}
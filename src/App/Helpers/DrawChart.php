<?php
namespace App\Helpers;

use PhpOffice\PhpSpreadsheet\Chart\Chart;
use PhpOffice\PhpSpreadsheet\Chart\DataSeries;
use PhpOffice\PhpSpreadsheet\Chart\DataSeriesValues;
use PhpOffice\PhpSpreadsheet\Chart\Legend;
use PhpOffice\PhpSpreadsheet\Chart\PlotArea;
use PhpOffice\PhpSpreadsheet\Chart\Title;
use PhpOffice\PhpSpreadsheet\Spreadsheet;
use PhpOffice\PhpSpreadsheet\Writer\Xlsx;
use PhpOffice\PhpSpreadsheet\IOFactory;
use PhpOffice\PhpSpreadsheet\Worksheet\Worksheet;


class DrawChart
{
    
    


    /**
     * The labels need to be the first column, skiping the headers
     * - this would then split charts by project, env, etc
     * - TEST!$A$2, TEST!$A$3 ets
     * - if we were using column headers it would be more like one TEST!$B$1, TEST!$C$1
     */
    protected static function dataSeriesLabels(Worksheet $worksheet, $data)
    {
        $dataSeriesLabels = [];
        /**
         * This works out the label name as a series over the row for merged data
         * such as ProjectEnvironment grouping
         */
        $alphabet = range('A', 'Z');
        $firstColumn = "A";
        $column = self::firstColumn($data, true);
        // there is no label column - eg, just totals
        if($column == 0) return $dataSeriesLabels;

        $last =  $column - 1;
        $lastColumn = $alphabet[$last];
        $worksheetName = $worksheet->getTitle();
        // there is always headers
        $firstRow = 2;
        $lastRowNumber = $worksheet->getHighestRow();
        
        for($row = $firstRow; $row <= $lastRowNumber; $row ++)
        {
            $ref = sprintf('%s!$%s$%s:$%s$%s', 
                        $worksheetName, 
                        $firstColumn,
                        $row,
                        $lastColumn,
                        $row                        
                    );
            
            $dataSeriesLabels[] = new DataSeriesValues(DataSeriesValues::DATASERIES_TYPE_STRING, $ref, null, 1);
        }
        return $dataSeriesLabels;
    }


    /**
     * Return the first column letter that contains a number
     */
    protected static function firstColumn($data, bool $returnInt = false) 
    {
        $alphabet = range('A', 'Z');
        $firstRow = $data[0];
        
        $i = 0;
        foreach($firstRow as $v)
        {
            if(is_numeric($v)) return ($returnInt) ? $i : $alphabet[$i];
            $i ++;
        }
        return null;
    }
    /**
     * We use the column headers as the x values - so first the first numeric col
     * - 'TEST!$B$1:$D$1'
     * If this was based on firt item in row, then it would be TEST!$A$2:$A$5
     */
    protected static function xAxis(Worksheet $worksheet, $data)
    {
        $lastColumn = $worksheet->getHighestColumn();
        $firstColumn = self::firstColumn($data);
        
        //'TEST!$B$1:$D$1'
        $ref = sprintf('%s!$%s$1:$%s$1',
            $worksheet->getTitle(),
            $firstColumn,
            $lastColumn
        );

        return [
            new DataSeriesValues(DataSeriesValues::DATASERIES_TYPE_STRING, $ref, null, 4), // Q1 to Q4
        ];

    }

    /**
     * This maps the row data, so starts at first numeric value to last column
     * per row
     */
    protected static function dataSeriesValues(Worksheet $worksheet, $data)
    {
        $lastColumn = $worksheet->getHighestColumn();
        $firstColumn = self::firstColumn($data);
        $rowCount = $worksheet->getHighestRow();
        $name = $worksheet->getTitle();
        $dataSeriesValues = [];
        for($row = 2; $row <= $rowCount; $row ++)
        {
            //'TEST!$B$2:$D$2'
            $ref = sprintf('%s!$%s$%s:$%s$%s', 
                        $name,
                        $firstColumn,
                        $row,
                        $lastColumn,
                        $row
                    );

            $dataSeriesValues[] = new DataSeriesValues(DataSeriesValues::DATASERIES_TYPE_NUMBER, $ref, null, 4);
        }
        return $dataSeriesValues;
    }
    
    /**
     * 
     */
    public static function draw(string $type, $data, string $sourceWorksheetName, string $targetWorksheetName, Spreadsheet &$spreadsheet, $stacked = DataSeries::GROUPING_STACKED) : void
    {
        $targetWorksheet = CostSpreadsheet::getWorksheet($spreadsheet, "Chart For ". $targetWorksheetName);
        $sourceWorksheet = CostSpreadsheet::getWorksheet($spreadsheet, $sourceWorksheetName);

        $dataSeriesLabels = self::dataSeriesLabels($sourceWorksheet, $data);
        $xAxisTickValues = self::xAxis($sourceWorksheet, $data);
        $dataSeriesValues = self::dataSeriesValues($sourceWorksheet, $data);
        
        // Build the dataseries
        $series = new DataSeries(
            $type,
            $stacked, // plotGrouping
            range(0, count($dataSeriesValues) - 1), // plotOrder
            $dataSeriesLabels, 
            $xAxisTickValues, 
            $dataSeriesValues        
        );
        // Set additional dataseries parameters
        //     Make it a horizontal bar rather than a vertical column graph
        $series->setPlotDirection(DataSeries::DIRECTION_COLUMN);
        
        // Set the series in the plot area
        $plotArea = new PlotArea(null, [$series]);
        // Set the chart legend
        $legend = new Legend(Legend::POSITION_RIGHT, null, false);
        
        $title = new Title('');
        $yAxisLabel = new Title('Costs ($)');
        $xAxisLabel = new Title('Year-Month');
        
        // Create the chart
        $chart = new Chart(
            'chart1', // name
            $title, // title
            $legend, // legend
            $plotArea, // plotArea
            true, // plotVisibleOnly
            DataSeries::EMPTY_AS_GAP, // displayBlanksAs
            $xAxisLabel, // xAxisLabel
            $yAxisLabel  // yAxisLabel
        );
        
        // Set the position where the chart should appear in the worksheet
        $chart->setTopLeftPosition('A1');
        $chart->setBottomRightPosition('AH59');
        
        // Add the chart to the worksheet
        $targetWorksheet->addChart($chart);
        
    }


    


}
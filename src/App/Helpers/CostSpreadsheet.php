<?php
namespace App\Helpers;

use PhpOffice\PhpSpreadsheet\Spreadsheet;
use PhpOffice\PhpSpreadsheet\Writer\Xlsx;
use PhpOffice\PhpSpreadsheet\IOFactory;
use PhpOffice\PhpSpreadsheet\Worksheet\Worksheet;


class CostSpreadsheet
{
    /**
     * Add data passed to a named worksheet
     * - doesnt save!
     */
    public static function writeToWorksheet($data, string $worksheetName, Spreadsheet &$spreadsheet, bool $hidesheet = false)
    {
        
        $worksheet = self::getWorksheet($spreadsheet, $worksheetName);
        
        if($hidesheet) $worksheet->setSheetState(Worksheet::SHEETSTATE_HIDDEN);

        $row = $worksheet->getHighestRow();  
        // add headings in
        if($row === 1) array_unshift( $data, array_keys($data[0]) );
        // write data from $row as worksheet is based on project name, will be many environments
        $worksheet->fromArray(
            $data,
            null,
            'A'.$row
        );
    }


    /**
     * Return existing spreadsheet or a new instance
     */
    public static function getSpreadsheet($filename): Spreadsheet
    {
        if(is_file($filename)) return IOFactory::load($filename);
        else return new Spreadsheet();
    }

    /**
     * Find named worksheet
     * - if not there, look for a default named worksheet & convert it
     * - otherwise, add a new worksheet and give it a name
     */
    public static function getWorksheet(Spreadsheet $spreadsheet, string $name)
    {
        $worksheets = self::getWorksheets($spreadsheet);
        if( in_array($name, $worksheets)) return $spreadsheet->getSheetByName($name);
        else if(in_array("Worksheet", $worksheets)) return $spreadsheet->getSheetByName("Worksheet")->setTitle($name);
        else return $spreadsheet->createSheet(0)->setTitle($name);
    }

    /**
     * Find all the worksheets
     */
    public static function getWorksheets(Spreadsheet $spreadsheet)
    {
        if($spreadsheet->getSheetCount() > 0) return $spreadsheet->getSheetNames();
        else return [];
    }

    




}
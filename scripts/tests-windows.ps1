# tests-windows.ps1
# Colors for output
$GREEN = 'Green'
$YELLOW = 'Yellow'
$RED = 'Red'
$NC = 'Gray'

# Configuration
$API_HOST = "http://localhost:8080"
$DB_HOST = "localhost"
$DB_PORT = "5432"
$DB_NAME = "project-sem-1"
$DB_USER = "validator"
$DB_PASSWORD = "val1dat0r"

# Получаем абсолютный путь к директории скрипта
$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path

# Temporary files for testing (используем абсолютные пути)
$TEST_CSV = Join-Path $SCRIPT_DIR "test_data.csv"
$TEST_ZIP = Join-Path $SCRIPT_DIR "test_data.zip"
$TEST_TAR = Join-Path $SCRIPT_DIR "test_data.tar"
$RESPONSE_ZIP = Join-Path $SCRIPT_DIR "response.zip"

function Create-TestFiles {
    param(
        [int]$Level
    )
    
    # Сохраняем текущую директорию
    $currentDir = Get-Location
    
    try {
        # Переходим в директорию скрипта для создания файлов
        Set-Location $SCRIPT_DIR
        
        if ($Level -eq 3) {
            # Create test CSV file with invalid data for complex level
            @"
id,name,category,price,create_date
1,item1,cat1,100,2024-01-01
2,item2,cat2,200,2024-01-15
3,item3,cat3,invalid_price,2024-01-20
4,,cat4,400,2024-01-25
5,item5,,500,2024-01-30
6,item6,cat6,600,invalid_date
1,item1,cat1,100,2024-01-01
"@ | Out-File -FilePath $TEST_CSV -Encoding UTF8
            
            if (Test-Path $TEST_ZIP) { Remove-Item $TEST_ZIP -Force }
            if (Test-Path $TEST_TAR) { Remove-Item $TEST_TAR -Force }
            
            Compress-Archive -Path $TEST_CSV -DestinationPath $TEST_ZIP -Force
            
            # Используем полный путь для tar
            $csvFileName = Split-Path $TEST_CSV -Leaf
            & tar -cf $TEST_TAR -C $SCRIPT_DIR $csvFileName
        }
        else {
            # Create test CSV file with valid data for simple and advanced levels
            @"
id,name,category,price,create_date
1,item1,cat1,100,2024-01-01
2,item2,cat2,200,2024-01-15
3,item3,cat3,300,2024-01-20
"@ | Out-File -FilePath $TEST_CSV -Encoding UTF8
            
            if (Test-Path $TEST_ZIP) { Remove-Item $TEST_ZIP -Force }
            if (Test-Path $TEST_TAR) { Remove-Item $TEST_TAR -Force }
            
            Compress-Archive -Path $TEST_CSV -DestinationPath $TEST_ZIP -Force
            
            # Используем полный путь для tar
            $csvFileName = Split-Path $TEST_CSV -Leaf
            & tar -cf $TEST_TAR -C $SCRIPT_DIR $csvFileName
        }
    }
    finally {
        # Возвращаемся в исходную директорию
        Set-Location $currentDir
    }
}

function Invoke-UploadRequest {
    param(
        [string]$Uri,
        [string]$FilePath,
        [string]$FormFieldName = "file"
    )
    
    # Проверяем существование файла
    if (-not (Test-Path $FilePath)) {
        Write-Host "[FAIL] File not found: $FilePath" -ForegroundColor $RED
        throw "File not found: $FilePath"
    }
    
    # Create multipart/form-data manually for PowerShell 5.1 compatibility
    $boundary = [System.Guid]::NewGuid().ToString()
    $fileBytes = [System.IO.File]::ReadAllBytes($FilePath)
    $fileName = [System.IO.Path]::GetFileName($FilePath)
    
    $bodyLines = @(
        "--$boundary",
        "Content-Disposition: form-data; name=`"$FormFieldName`"; filename=`"$fileName`"",
        "Content-Type: application/octet-stream",
        "",
        [System.Text.Encoding]::GetEncoding('iso-8859-1').GetString($fileBytes),
        "--$boundary--",
        ""
    )
    
    $body = $bodyLines -join "`r`n"
    
    try {
        $response = Invoke-RestMethod -Uri $Uri `
            -Method Post `
            -ContentType "multipart/form-data; boundary=$boundary" `
            -Body $body
        
        return $response
    }
    catch {
        Write-Host "[FAIL] Request failed: $($_.Exception.Message)" -ForegroundColor $RED
        throw
    }
}

function Check-ApiSimple {
    Create-TestFiles -Level 1

    Write-Host "`nAPI Check (Simple Level)" -ForegroundColor $YELLOW
    
    # Check POST /api/v0/prices
    Write-Host "Testing POST /api/v0/prices" -ForegroundColor $YELLOW
    
    try {
        $response = Invoke-UploadRequest -Uri "${API_HOST}/api/v0/prices" -FilePath $TEST_ZIP
        
        # Проверяем наличие всех трёх полей в ответе (аналогично bash)
        if ($response.PSObject.Properties.Name -contains 'total_items' -and 
            $response.PSObject.Properties.Name -contains 'total_categories' -and 
            $response.PSObject.Properties.Name -contains 'total_price') {
            Write-Host "[OK] POST request successful" -ForegroundColor $GREEN
        }
        else {
            Write-Host "[FAIL] POST request failed" -ForegroundColor $RED
            return $false
        }
    }
    catch {
        Write-Host "[FAIL] POST request failed: $_" -ForegroundColor $RED
        return $false
    }
    
    # Check GET /api/v0/prices
    Write-Host "Testing GET /api/v0/prices" -ForegroundColor $YELLOW
    
    $currentDir = Get-Location
    $tmpDir = Join-Path $env:TEMP ([System.Guid]::NewGuid().ToString())
    New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null
    
    try {
        Set-Location $tmpDir
        
        Invoke-WebRequest -Uri "${API_HOST}/api/v0/prices" -OutFile $RESPONSE_ZIP
        
        # Unzip archive
        Expand-Archive -Path $RESPONSE_ZIP -DestinationPath . -Force
        
        if (Test-Path "data.csv") {
            Write-Host "[OK] GET request successful" -ForegroundColor $GREEN
            Set-Location $currentDir
            Remove-Item $tmpDir -Recurse -Force
            return $true
        }
        else {
            Write-Host "[FAIL] GET request returned invalid archive" -ForegroundColor $RED
            Set-Location $currentDir
            Remove-Item $tmpDir -Recurse -Force
            return $false
        }
    }
    catch {
        Write-Host "[FAIL] GET request failed: $_" -ForegroundColor $RED
        Set-Location $currentDir
        if (Test-Path $tmpDir) {
            Remove-Item $tmpDir -Recurse -Force
        }
        return $false
    }
}

function Check-ApiAdvanced {
    Create-TestFiles -Level 2
    
    Write-Host "`nAPI Check (Advanced Level)" -ForegroundColor $YELLOW
    
    # Check POST with ZIP
    Write-Host "Testing POST /api/v0/prices?type=zip" -ForegroundColor $YELLOW
    
    try {
        $response = Invoke-UploadRequest -Uri "${API_HOST}/api/v0/prices?type=zip" -FilePath $TEST_ZIP
        
        if ($response.PSObject.Properties.Name -contains 'total_items') {
            Write-Host "[OK] POST request with ZIP successful" -ForegroundColor $GREEN
        }
        else {
            Write-Host "[FAIL] POST request with ZIP failed" -ForegroundColor $RED
            return $false
        }
    }
    catch {
        Write-Host "[FAIL] POST request with ZIP failed: $_" -ForegroundColor $RED
        return $false
    }
    
    # Check POST with TAR
    Write-Host "Testing POST /api/v0/prices?type=tar" -ForegroundColor $YELLOW
    
    try {
        $response = Invoke-UploadRequest -Uri "${API_HOST}/api/v0/prices?type=tar" -FilePath $TEST_TAR
        
        if ($response.PSObject.Properties.Name -contains 'total_items') {
            Write-Host "[OK] POST request with TAR successful" -ForegroundColor $GREEN
        }
        else {
            Write-Host "[FAIL] POST request with TAR failed" -ForegroundColor $RED
            return $false
        }
    }
    catch {
        Write-Host "[FAIL] POST request with TAR failed: $_" -ForegroundColor $RED
        return $false
    }
    
    # Check GET
    return Check-ApiSimple
}

function Check-ApiComplex {
    Create-TestFiles -Level 3
    Write-Host "`nAPI Check (Complex Level)" -ForegroundColor $YELLOW
    
    # Check POST with invalid data processing
    Write-Host "Testing POST /api/v0/prices?type=zip with invalid data" -ForegroundColor $YELLOW
    
    try {
        $response = Invoke-UploadRequest -Uri "${API_HOST}/api/v0/prices?type=zip" -FilePath $TEST_ZIP
        
        # Check all required fields in response (теперь проверяем наличие через -contains)
        $required_fields = @("total_count", "duplicates_count", "total_items", "total_categories", "total_price")
        $missing_fields = @()
        
        foreach ($field in $required_fields) {
            if (-not ($response.PSObject.Properties.Name -contains $field)) {
                $missing_fields += $field
            }
        }
        
        if ($missing_fields.Count -eq 0) {
            Write-Host "[OK] All required fields present in response" -ForegroundColor $GREEN
        }
        else {
            Write-Host "[FAIL] Missing required fields: $($missing_fields -join ', ')" -ForegroundColor $RED
            return $false
        }
        
        # Check correct handling of invalid data (эти проверки оставляем как есть, т.к. они работают с числовыми значениями)
        if ($response.total_count -gt $response.total_items) {
            Write-Host "[OK] Invalid data correctly filtered (total_count > total_items)" -ForegroundColor $GREEN
        }
        else {
            Write-Host "[FAIL] Problem with invalid data processing" -ForegroundColor $RED
            return $false
        }
        
        # Check duplicate handling (эта проверка оставлена как есть, т.к. работает с числом)
        if ($response.duplicates_count -gt 0) {
            Write-Host "[OK] Duplicates successfully detected" -ForegroundColor $GREEN
        }
        else {
            Write-Host "[FAIL] Problem with duplicate detection" -ForegroundColor $RED
            return $false
        }
    }
    catch {
        Write-Host "[FAIL] POST request failed: $_" -ForegroundColor $RED
        return $false
    }
    
    # Check GET with filters
    Write-Host "Testing GET /api/v0/prices with filters" -ForegroundColor $YELLOW
    $filters = "start=2024-01-01&end=2024-01-31&min=30&max=1000"
    
    $currentDir = Get-Location
    $tmpDir = Join-Path $env:TEMP ([System.Guid]::NewGuid().ToString())
    New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null
    
    try {
        Set-Location $tmpDir
        
        Invoke-WebRequest -Uri "${API_HOST}/api/v0/prices?${filters}" -OutFile $RESPONSE_ZIP
        
        # Unzip archive
        Expand-Archive -Path $RESPONSE_ZIP -DestinationPath . -Force
        
        if (-not (Test-Path "data.csv")) {
            Write-Host "[FAIL] data.csv file not found in archive" -ForegroundColor $RED
            Set-Location $currentDir
            Remove-Item $tmpDir -Recurse -Force
            return $false
        }
        
        # Check that exported data doesn't contain invalid records
        $data = Get-Content "data.csv" -Raw
        $invalidPattern = @(
            ',invalid_',
            '^[^,]*,[^,]*,,[^,]*,[^,]*',
            '^[^,]*,,[^,]*,[^,]*,[^,]*'
        )
        
        $hasInvalid = $false
        foreach ($pattern in $invalidPattern) {
            if ($data -match $pattern) {
                $hasInvalid = $true
                break
            }
        }
        
        if (-not $hasInvalid) {
            Write-Host "[OK] Exported data doesn't contain invalid records" -ForegroundColor $GREEN
        }
        else {
            Write-Host "[FAIL] Invalid records found in export" -ForegroundColor $RED
            Set-Location $currentDir
            Remove-Item $tmpDir -Recurse -Force
            return $false
        }
        
        Set-Location $currentDir
        Remove-Item $tmpDir -Recurse -Force
        return $true
    }
    catch {
        Write-Host "[FAIL] GET request with filters failed: $_" -ForegroundColor $RED
        Set-Location $currentDir
        if (Test-Path $tmpDir) {
            Remove-Item $tmpDir -Recurse -Force
        }
        return $false
    }
}

function Check-Postgres {
    param(
        [int]$Level
    )
    
    Write-Host "`nPostgreSQL Check (Level $Level)" -ForegroundColor $YELLOW
    
    # Basic connection check for all levels
    try {
        $env:PGPASSWORD = $DB_PASSWORD
        & psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c '\q' 2>$null
    }
    catch {
        Write-Host "[FAIL] PostgreSQL unavailable" -ForegroundColor $RED
        return $false
    }
    
    switch ($Level) {
        1 {
            Write-Host "Executing level 1 check" -ForegroundColor $YELLOW
            try {
                $query = "SELECT COUNT(*) FROM prices;"
                $env:PGPASSWORD = $DB_PASSWORD
                $result = & psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c $query 2>$null
                
                if ($LASTEXITCODE -eq 0) {
                    Write-Host "[OK] PostgreSQL working correctly" -ForegroundColor $GREEN
                    return $true
                }
                else {
                    Write-Host "[FAIL] Query execution error" -ForegroundColor $RED
                    return $false
                }
            }
            catch {
                Write-Host "[FAIL] Query execution error: $_" -ForegroundColor $RED
                return $false
            }
        }
        
        2 {
            Write-Host "Executing level 2 check" -ForegroundColor $YELLOW
            try {
                $query = @"
                SELECT 
                    COUNT(*) as total_items,
                    COUNT(DISTINCT category) as total_categories,
                    SUM(price) as total_price
                FROM prices;
"@
                $env:PGPASSWORD = $DB_PASSWORD
                $result = & psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c $query 2>$null
                
                if ($LASTEXITCODE -eq 0) {
                    Write-Host "[OK] PostgreSQL working correctly" -ForegroundColor $GREEN
                    return $true
                }
                else {
                    Write-Host "[FAIL] Query execution error" -ForegroundColor $RED
                    return $false
                }
            }
            catch {
                Write-Host "[FAIL] Query execution error: $_" -ForegroundColor $RED
                return $false
            }
        }
        
        3 {
            Write-Host "Executing level 3 check" -ForegroundColor $YELLOW
            try {
                $query = @"
                WITH stats AS (
                    SELECT 
                        COUNT(*) as total_items,
                        COUNT(DISTINCT category) as total_categories,
                        SUM(price) as total_price,
                        COUNT(*) - COUNT(DISTINCT (name, category, price)) as duplicates
                    FROM prices
                    WHERE create_date BETWEEN '2024-01-01' AND '2024-01-31'
                    AND price BETWEEN 300 AND 1000
                )
                SELECT * FROM stats;
"@
                $env:PGPASSWORD = $DB_PASSWORD
                # Убираем перенаправление stderr в $null, чтобы видеть ошибки
                $result = & psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c $query 2>&1
        
                
                if ($LASTEXITCODE -eq 0) {
                    Write-Host "[OK] PostgreSQL working correctly" -ForegroundColor $GREEN
                    return $true
                }
                else {
                    Write-Host "[FAIL] Query execution error" -ForegroundColor $RED
                    return $false
                }
            }
            catch {
                Write-Host "[FAIL] Query execution error: $_" -ForegroundColor $RED
                return $false
            }
        }
        
        default {
            Write-Host "[FAIL] Unknown level: $Level" -ForegroundColor $RED
            return $false
        }
    }
}

function Cleanup {
    # Удаляем временные файлы в директории скрипта
    $files = @($TEST_CSV, $TEST_ZIP, $TEST_TAR, $RESPONSE_ZIP)
    foreach ($file in $files) {
        if (Test-Path $file) {
            Remove-Item -Path $file -ErrorAction SilentlyContinue
        }
    }
}

function Main {
    param(
        [int]$Level
    )
    
    $failed = 0
    
    switch ($Level) {
        1 {
            Write-Host "=== Starting Simple Level Check ===" -ForegroundColor $YELLOW
            if (-not (Check-ApiSimple)) { $failed++ }
            if (-not (Check-Postgres -Level 1)) { $failed++ }
        }
        2 {
            Write-Host "=== Starting Advanced Level Check ===" -ForegroundColor $YELLOW
            if (-not (Check-ApiAdvanced)) { $failed++ }
            if (-not (Check-Postgres -Level 2)) { $failed++ }
        }
        3 {
            Write-Host "=== Starting Complex Level Check ===" -ForegroundColor $YELLOW
            if (-not (Check-ApiComplex)) { $failed++ }
            if (-not (Check-Postgres -Level 3)) { $failed++ }
        }
        default {
            Write-Host "[FAIL] Invalid check level" -ForegroundColor $RED
            Cleanup
            exit 1
        }
    }
    
    Cleanup
    
    Write-Host "`nCheck Summary:" -ForegroundColor $YELLOW
    if ($failed -eq 0) {
        Write-Host "[OK] All checks passed successfully" -ForegroundColor $GREEN
        exit 0
    }
    else {
        Write-Host "[FAIL] Problems found in $failed checks" -ForegroundColor $RED
        exit 1
    }
}

# Argument check
if ($args.Count -ne 1 -or $args[0] -notmatch '^[1-3]$') {
    Write-Host "Usage: $PSCommandPath <check_level>"
    Write-Host "Check level must be:"
    Write-Host "  1 - simple level"
    Write-Host "  2 - advanced level"
    Write-Host "  3 - complex level"
    exit 1
}

Main -Level $args[0]
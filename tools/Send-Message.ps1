function GetRandomString {
    $length = 8
    $characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    
    $random = New-Object System.Random
    $randomString = -join (0..($length - 1) | ForEach-Object { $characters[$random.Next(0, $characters.Length)] })
    
    return $randomString
}

function EnqueueRandomMessages($url, $count) {
    for ($i = 0; $i -lt $count; $i++) {
        $message = GetRandomString
        $body = @{ message = $message } | ConvertTo-Json
        # Invoke-RestMethod -Uri $url -Method POST -Body $body -ContentType "application/json"

        try {
            Invoke-RestMethod -Uri $url -Method POST -Body $body -ContentType "application/json"
        } catch {
            # Dig into the exception to get the Response details.
            # Note that value__ is not a typo.
            Write-Host "StatusCode:" $_.Exception.Response.StatusCode.value__ 
            Write-Host "StatusDescription:" $_.Exception.Response.StatusDescription
        }

        Write-Host "Enqueued message: $message"
    }
}

$Url = "https://keda-producer.fmz-c-x-app-aks-01.corp.fmglobal.com/enqueue"
$Count = 100

EnqueueRandomMessages $Url $Count



Invoke-RestMethod -Uri $url -Method POST -Body $body -ContentType "application/json"
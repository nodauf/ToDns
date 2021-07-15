function DownloadOverDNS {
<#
.SYNOPSIS
Function to download a file through DNS

.DESCRIPTION
Will made TXT request for <id>.d.<domain> until the dn server returned an error. The base64 payload will be decoded and wrote to a file.

.PARAMETER Domain
The domain where the todns script listen to be able to answer the requests.

.PARAMETER OutputFilename
Filename to write the base64 decoded data.

.EXAMPLE
PS > DownloadOverDNS -Domain todns.nodauf.ovh -OutputFilename out.exe

.LINK
https://github.com/nodauf/ToDns
#>
    param(
        [Parameter(Mandatory=$true)]$Domain,
        [Parameter(Mandatory=$true)]$OutputFilename
    )
    $finish = $false
    $i = 0
    $dataBase64 = ""
    do {
        try{
            $dataBase64 += (Resolve-DnsName "$i.d.$Domain" -Type TXT  -ErrorAction Stop).Strings
            Write-Host "Download chunk $i at $i.d.$Domain"
            $i++
        }catch{
            $finish = $true
        }
    } while (!$finish)
   $dataBinary = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($dataBase64))
   $decoded = [System.Convert]::FromBase64CharArray($dataBase64, 0, $dataBase64.Length)

   Set-Content $OutputFilename -Value $decoded -Encoding Byte
}

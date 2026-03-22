echo Allowing Firewall Rules
powershell -Command "Start-Process powershell -Verb RunAs -ArgumentList '-File \"%~dp0allow-firewall-rules.ps1\"'"  

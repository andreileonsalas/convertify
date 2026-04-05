; Convertify - Inno Setup Script
; Download Inno Setup from: https://jrsoftware.org/isinfo.php

#define AppName "Convertify"
#define AppVersion "1.0.0"
#define AppPublisher "Convertify"
#define AppURL "https://github.com/youruser/convertify"
#define AppExeName "convertify.exe"

[Setup]
AppId={{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}
AppName={#AppName}
AppVersion={#AppVersion}
AppPublisher={#AppPublisher}
AppPublisherURL={#AppURL}
AppSupportURL={#AppURL}
AppUpdatesURL={#AppURL}
DefaultDirName={autopf}\{#AppName}
DefaultGroupName={#AppName}
AllowNoIcons=yes
OutputDir=dist
OutputBaseFilename=convertify-setup
Compression=lzma
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=admin

[Languages]
Name: "spanish"; MessagesFile: "compiler:Languages\Spanish.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
; Main executable
Source: "dist\{#AppExeName}"; DestDir: "{app}"; Flags: ignoreversion

; FFmpeg binaries (place ffmpeg.exe and ffprobe.exe in assets\ before building)
Source: "assets\ffmpeg.exe";  DestDir: "{app}"; Flags: ignoreversion
Source: "assets\ffprobe.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\{#AppName}";                    Filename: "{app}\{#AppExeName}"
Name: "{group}\Desinstalar {#AppName}";        Filename: "{uninstallexe}"
Name: "{commondesktop}\{#AppName}";            Filename: "{app}\{#AppExeName}"; Tasks: desktopicon

[Run]
; After install, register context menus automatically
Filename: "{app}\{#AppExeName}"; Parameters: "install"; \
  Description: "Registrar menús contextuales"; \
  Flags: runhidden waituntilterminated

; Optionally launch the app after install
Filename: "{app}\{#AppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(AppName, '&', '&&')}}"; \
  Flags: nowait postinstall skipifsilent

[UninstallRun]
; On uninstall, remove context menus first
Filename: "{app}\{#AppExeName}"; Parameters: "uninstall"; Flags: runhidden waituntilterminated

[Code]
// Optional: check if FFmpeg files exist before install
function InitializeSetup(): Boolean;
begin
  Result := True;
end;

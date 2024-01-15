{ pkgs, ... }:

{
  packages = with pkgs; [
    air
  ];

  dotenv.disableHint = true;
  languages.go.enable = true;
  services.mongodb.enable = true;

  scripts.dev.exec = ''
    ${pkgs.air}/bin/air
  '';
}
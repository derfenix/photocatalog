{
  description = "Photo/video organization tool";

  inputs.nixpkgs.url = "nixpkgs/nixos-24.11";

  outputs = { self, nixpkgs }:
    let
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = "2.0.0";
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          photocatalog = pkgs.buildGoModule {
            pname = "photocatalog";
            inherit version;
            src = ./.;
            vendorHash = "sha256-dj11SRRoB8ZbkcQs75HPI0DpW4c5jzY0N8MD1wKpw+4=";
          };
        }
      );

      nixosModules.photocatalog = { config, lib, pkgs, ... }:
        with lib;
      {
        options.photocatalog = {
          enable = lib.mkEnableOption "Enable photocatalog";

          syncs = mkOption {
            default = {};
            description = ''
              Organization paths with its own params.
            '';
            example = {

            };
            type = types.attrsOf (types.submodule ({ name, ... }: {
              freeformType = settingsFormat.type;
              options = {
                source = mkOption {
                  type = types.str;
                  default = name;
                  description = ''
                    Source folder path.
                  '';
                };
                target = mkOption {
                  type = types.str;
                  description = ''
                    Target folder path.
                  '';
                };
                overwrite = mkOption {
                  type = types.bool;
                  default = false;
                  description = ''
                    Overwrite files, existing in target.
                  '';
                };
                watch = mkOption {
                  type = types.bool;
                  default = true;
                  description = ''
                    Watch for new files in source path.
                  '';
                };
                skipFullSync = mkOption {
                  type = types.bool;
                  default = false;
                  description = ''
                    Do not make full sync.
                  '';
                };
              };
            }));
          };
        };

        config = lib.mkIf config.photocatalog.enable {
          environment.systemPackages = [ self.packages.${pkgs.system}.photocatalog ];
          systemd.services = lib.genAttrs config.photocatalog.syncs (sync:
            {
              ${sync.name} = {
                name = "photocatalog_${sync.name}";
                after = [ "local-fs.target" ];
                path = [
                  self.packages.${pkgs.system}.photocatalog
                ];
                preStart = lib.mkIf (!sync.skipFullSync) [
                  "mkdir -p ${sync.target}"
                  "photocatalog -source ${sync.source} -target ${sync.target}"
                ];
                script = [
                  "photocatalog -source ${sync.source} -target ${sync.target} -skip-full-sync -watch"
                ];
              };
            }
          );
        };
      };

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools ];
          };
        });
      defaultPackage = forAllSystems (system: self.packages.${system}.photocatalog);
    };
}

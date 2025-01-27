{
  description = "Photo/video organization tool";

  inputs.nixpkgs.url = "nixpkgs/nixos-24.11";

  outputs = { self, nixpkgs }:
    let
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = "2.0.4";
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
                mode = mkOption {
                  type = types.enum [ "hardlink" "symlink" "move" "copy" ];
                  default = "hardlink";
                  description = ''
                    Organization mode, one of [ hardlink symlink move copy ].
                  '';
                };
              };
            }));
          };
        };

        config = lib.mkIf config.photocatalog.enable {
          environment.systemPackages = [ self.packages.${pkgs.system}.photocatalog ];
          systemd.user.services = lib.mapAttrs' (name: sync: nameValuePair
            ("photocatalog${lib.replaceStrings ["/"] ["-"] sync.source}")
            {
                after = [ "local-fs.target" ];
                path = [
                  self.packages.${pkgs.system}.photocatalog
                ];
                wantedBy = [
                  "default.target"
                ];
                preStart = ''
                  mkdir -p ${sync.target}
                '' + (if !sync.skipFullSync then (''
                  photocatalog -source ${sync.source} -target ${sync.target} -mode ${sync.mode} ${if sync.overwrite then "-overwrite" else ""}
                '') else null);
                script = "photocatalog -source ${sync.source} -target ${sync.target} -skip-full-sync -watch -mode ${sync.mode} ${if sync.overwrite then "-overwrite" else ""}";
                serviceConfig = {
                  Type="simple";
                  Restart="no";
                };
            }
          ) config.photocatalog.syncs;
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

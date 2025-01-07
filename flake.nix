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
        });

    nixosModules.photocatalog = { config, lib, pkgs, ... }: {
          options.photocatalog = {
            enable = lib.mkEnableOption "Enable photocatalog";
          };

          config = lib.mkIf config.photocatalog.enable {
            environment.systemPackages = [ self.packages.${pkgs.system}.photocatalog ];
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

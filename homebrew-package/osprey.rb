class Osprey < Formula
  desc "Modern functional programming language designed for clarity, safety, and expressiveness"
  homepage "https://www.ospreylang.dev"
  url "https://github.com/melbournedeveloper/osprey/releases/download/v0.1.0/osprey-darwin-amd64.tar.gz"
  version "0.1.0"
  sha256 "PLACEHOLDER_SHA256"
  
  depends_on "llvm"

  def install
    # Install pre-built binaries and libraries
    bin.install "osprey"
    lib.install "libfiber_runtime.a"
    lib.install "libhttp_runtime.a"
  end

  test do
    # Test that the compiler can show help
    output = shell_output("#{bin}/osprey --help 2>&1", 0)
    assert_match "Osprey", output
    
    # Test that runtime libraries are installed
    assert_predicate lib/"libfiber_runtime.a", :exist?
    assert_predicate lib/"libhttp_runtime.a", :exist?
  end
end 
Output ./docs/images/journey.gif

Set Shell "bash"
Set FontSize 18
Set FontFamily "FiraCode Nerd Font Mono"
Set Width 1600
Set Height 1200

# Launch lazypkg
Type "./lazypkg --demo" Sleep 500ms Enter

Sleep 2s Enter Sleep 1s

# -------------------------------------
# Package list
# -------------------------------------

# Update the second package
Down Sleep 1500ms
Type "u" Sleep 2s Enter
Sleep 3s

# Back to the sidebar
Type "h" Sleep 2s

# -------------------------------------
# Sidebar
# -------------------------------------

# Select the second package manager
Sleep 500ms Type "j" Sleep 1500ms
Enter Sleep 2s

# Update the all packages
Type "a" Sleep 3s Enter

Sleep 8s

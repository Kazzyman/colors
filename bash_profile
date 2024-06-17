# To implement gls command.
eval "$(/opt/homebrew/bin/brew shellenv)"

# Reads the .dircolors configuration file used by gls command. 
eval "$(gdircolors -b ~/.dircolors)"

# Also relevent to the gls command, gdircolors CAN optionally create an initial .dircolors file. Only run this once!
# ... as it will actually clobber any existing (edited) .dircolors file if run post-edits to said file. 
# gdircolors -p > ~/.dircolors

# This next section is relevent to the native ls command.
# ~/.bashrc contains: export LSCOLORS="exfxcxdxbxegedabagacad"
# older lscolors:
#export LSCOLORS=gxfxcxdxbxegedabagacad
# newer and favored colors:
export LSCOLORS=gxfxbxdxcxegedabagacad
# Ex,Fx,Bx,Dx,Cx,eg,ed,ab,ag,ac,ad
# in the above, there are 11 pairs, e.g., Ex and Fx are the first two pairs (bold blue, and bold magenta, both with default BG)
# x means default color (often black) and as the second member of a pair it refers to the background field
# : a: Black, b: Red, c: Green, d: Brown, e: Blue, f: Magenta, g: Cyan, h: Light gray
# : A: Bold black, usually shows up as dark gray
# ... B: Bold red, C: Bold green, D: Bold brown, usually shows up as yellow
# : E: Bold blue, F: Bold magenta, G: Bold cyan
# H: Bold light gray; looks like bright white when used as foreground color


# Format the prompt, and turn-on color. 
export PS1="\[\033[36m\]\u\[\033[m\]@\[\033[32m\]\h:\[\033[33;1m\]\w\[\033[m\]\$ "
export CLICOLOR=1


# Instead of using mundane and limited direct Aliasing for ls and gls, use these scrips. 
alias l='~/ricks.l.sh'
alias g='~/ricks.g.sh'
alias dir='~/ricks.g.sh'

alias cl=clear

# To set-up for using my song/file counting program.
alias cdm='cd ~/Music/Music/Media.localized/Music'

# a nice-enough python script for pretty-printing the output from gls:
alias gr='gls -oFGtpsrh --color=auto --time-style=long-iso --block-size=1 --group-directories-first | ~/ricks.colors.and.commas.py'

# my MOST-FAVORITE go program for pretty-printing the output from gls:
alias fl='gls -oFGtpsrh --color=auto --time-style=full-iso --block-size=1 --group-directories-first | ~/ricks.colors.and.commas'
alias fld='gls -oFGtpsrh --color=auto --time-style=full-iso --block-size=1 --group-directories-first | ~/ricks.colors.and.commas.DirOnly'

# gls works similarly to ls but was missing the , option and the T option (and I had to have commas! And, endless colors!)

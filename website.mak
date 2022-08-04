#
# Makefile for running pandoc on all Markdown docs ending in .md
#
PROJECT = shorthand

MD_PAGES = $(shell ls -1 *.md | grep -v "nav.md")

HTML_PAGES = $(shell ls -1 *.md | grep -v "nav.md" | sed -E 's/.md/.html/g')

build: $(HTML_PAGES) $(MD_PAGES) license.html shorthand.html shorthand-syntax.html shorthand-tutorial.html

$(HTML_PAGES): $(MD_PAGES) .FORCE
	pandoc --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html \
	    --template=page.tmpl
	@if [ $@ = "README.html" ]; then mv README.html index.html; fi

license.html: LICENSE
	pandoc --metadata title="$(PROJECT): License" -s --from Markdown --to html5 LICENSE -o license.html \
	    --template=page.tmpl

shorthand.html: docs/shorthand.md
	pandoc --metadata title="$(PROJECT): interpreter" -s --from Markdown --to html5 docs/shorthand.md -o shorthand.html --template page.tmpl
	
shorthand-syntax.html: docs/shorthand-syntax.md
	pandoc --metadata title="$(PROJECT): syntax" -s --from Markdown --to html5 docs/shorthand-syntax.md -o shorthand-syntax.html --template page.tmpl

shorthand-tutorial.html: docs/shorthand-tutorial.md
	pandoc --metadata title="$(PROJECT): tutorial" -s --from Markdown --to html5 docs/shorthand-tutorial.md -o shorthand-tutorial.html --template page.tmpl

clean:
	@if [ -f index.html ]; then rm *.html; fi

.FORCE:

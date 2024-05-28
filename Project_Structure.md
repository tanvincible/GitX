# Project Structure

Below is the comprehensive project structure, including all directory and file listings.

```
GitX
│   go.mod                   # Go module file
│   LICENSE                  # License file
│   README.md                # Readme file
│
├───cmd/                     # Main application entry point
│       main.go
│
├───docs/                    # Documentation files
│   │   404.md
│   │   Appendix.md
│   │   Changelog.md
│   │   Command_Reference.md
│   │   Contact.md
│   │   FAQS.md
│   │   feed.xml
│   │   Gemfile
│   │   Gemfile.lock
│   │   Getting_Started.md
│   │   index.markdown
│   │   LICENSE.markdown
│   │   Overview.md
│   │   search.json
│   │   sitemap.xml
│   │   tooltips.html
│   │   tooltips.json
│   │   Troubleshooting.md
│   │   Tutorials.md
│   │   User_Guide.md
│   │   _config.yml
│   │
│   ├───css/                 # CSS styles for documentation
│   │   │   bootstrap.min.css
│   │   │   boxshadowproperties.css
│   │   │   customstyles.css
│   │   │   font-awesome.min.css
│   │   │   modern-business.css
│   │   │   printstyles.css
│   │   │   syntax.css
│   │   │   theme-blue.css
│   │   │   theme-green.css
│   │   │
│   │   └───fonts/
│   │           fontawesome-webfont.eot
│   │           fontawesome-webfont.svg
│   │           fontawesome-webfont.ttf
│   │           fontawesome-webfont.woff
│   │           fontawesome-webfont.woff2
│   │           FontAwesome.otf
│   │
│   ├───fonts/               # Font files for documentation
│   │       fontawesome-webfont.eot
│   │       fontawesome-webfont.svg
│   │       fontawesome-webfont.ttf
│   │       fontawesome-webfont.woff
│   │       FontAwesome.otf
│   │       glyphicons-halflings-regular.eot
│   │       glyphicons-halflings-regular.svg
│   │       glyphicons-halflings-regular.ttf
│   │       glyphicons-halflings-regular.woff
│   │       glyphicons-halflings-regular.woff2
│   │
│   ├───js/                  # JavaScript files for documentation
│   │       customscripts.js
│   │       jekyll-search.js
│   │       jquery.ba-throttle-debounce.min.js
│   │       jquery.navgoco.min.js
│   │       jquery.shuffle.min.js
│   │       toc.js
│   │
│   ├───_data/               # Data files for documentation
│   │   │   alerts.yml
│   │   │   definitions.yml
│   │   │   glossary.yml
│   │   │   samplelist.yml
│   │   │   strings.yml
│   │   │   tags.yml
│   │   │   terms.yml
│   │   │   topnav.yml
│   │   │
│   │   └───sidebars
│   │           home_sidebar.yml
│   │
│   ├───_includes/           # Include files for documentation
│   │   │   archive.html
│   │   │   callout.html
│   │   │   feedback.html
│   │   │   footer.html
│   │   │   head.html
│   │   │   image.html
│   │   │   important.html
│   │   │   initialize_shuffle.html
│   │   │   inline_image.html
│   │   │   links.html
│   │   │   sidebar.html
│   │   │   taglogic.html
│   │   │   tip.html
│   │   │   toc.html
│   │   │   topnav.html
│   │   │   warning.html
│   │   │
│   │   └───custom
│   │           getting_started_series.html
│   │           getting_started_series_next.html
│   │           series_acme.html
│   │           series_acme_next.html
│   │           usermap.html
│   │           usermapcomplex.html
│   │
│   └───_layouts/            # Layout files for documentation
│           default.html
│           page.html
│
├───internal/                # Internal packages
│   │   merkletree.go
│   │
│   ├───compression/         # Compression logic
│   │       compression.go
│   │
│   ├───hash/                # Hashing logic
│   │       hash.go
│   │
│   ├───metadata/            # Metadata handling
│   │       metadata.go
│   │
│   └───storage/             # Storage handling
│           storage.go
│
├───models/                  # Data models
│       blob.go
│       commit.go
│       structs.go
│       tree.go
│
└───utils/                   # Utility packages
    ├───file_operations/     # File operations utility
    │       file_operations.go
    │
    ├───metadata_operations/ # Metadata operations utility
    │       metadata_operations.go
    │
    └───vcs_operations/      # Version control operations utility
            vcs_operations.go
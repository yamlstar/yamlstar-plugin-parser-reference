MAKES-COMMIT := a7b80ec8f10ac700693278a559f315fccd50b5ad
M ?= .cache/makes
$(shell test -d $M || { \
  git clone -q https://github.com/makeplus/makes $M && \
  git -C $M checkout -q $(MAKES-COMMIT); \
})

GO-VERSION := 1.27.1
include $M/init.mk
include $M/gloat.mk
include $M/go.mk
include $M/clean.mk
include $M/shell.mk

VERSION := 0.2.5
MODULE := github.com/yamlstar/yamlstar-plugin-parser-reference
YAML-PARSER-VERSION := 0.2.5
YAML-PARSER-FILE := yaml-parser-$(YAML-PARSER-VERSION).jar
YAML-PARSER-JAR := .cache/$(YAML-PARSER-FILE)
YAML-PARSER-BASE-URL := https://repo.clojars.org/org/yamlstar/yaml-parser
YAML-PARSER-URL := \
  $(YAML-PARSER-BASE-URL)/$(YAML-PARSER-VERSION)/$(YAML-PARSER-FILE)
YAML-PARSER-SRC-DIR := .cache/yaml-parser-$(YAML-PARSER-VERSION)
YAML-PARSER-SRC-STAMP := $(YAML-PARSER-SRC-DIR)/.extracted
SOURCE_CACHE := .cache/src
GENERATED_WORK := .cache/generated
GENERATED_DIR := internal/glojure

PARSER_SOURCES := \
  $(SOURCE_CACHE)/yaml_parser/prelude.clj \
  $(SOURCE_CACHE)/yaml_parser/parser.clj \
  $(SOURCE_CACHE)/yaml_parser/receiver.clj \
  $(SOURCE_CACHE)/yaml_parser/grammar.clj \
  $(SOURCE_CACHE)/yaml_parser/core.clj

MAKES-CLEAN := \
  $(YAML-PARSER-JAR) \
  $(YAML-PARSER-SRC-DIR) \
  $(SOURCE_CACHE) \
  $(GENERATED_WORK)

default:: test

test: $(GO)
	$(GO) test ./...

release-check: test check-generated
	grep -Fq 'const Version = "$(VERSION)"' parser/parser.go

generate: $(PARSER_SOURCES) $(GLOAT)
	rm -rf $(GENERATED_WORK)
	mkdir -p $(GENERATED_WORK)
	env -u GOROOT $(GLOAT) --force --module $(MODULE) \
	  -o $(GENERATED_WORK)/ $(PARSER_SOURCES)
	test -f $(GENERATED_WORK)/pkg/yaml_parser/core/loader.go
	rm -rf $(GENERATED_DIR)/pkg/yaml_parser
	mkdir -p $(GENERATED_DIR)/pkg
	cp -R $(GENERATED_WORK)/pkg/yaml_parser $(GENERATED_DIR)/pkg/

check-generated: $(PARSER_SOURCES) $(GLOAT)
	rm -rf $(GENERATED_WORK)
	mkdir -p $(GENERATED_WORK)
	env -u GOROOT $(GLOAT) --force --module $(MODULE) \
	  -o $(GENERATED_WORK)/ $(PARSER_SOURCES)
	diff -ru $(GENERATED_DIR)/pkg/yaml_parser \
	  $(GENERATED_WORK)/pkg/yaml_parser

$(YAML-PARSER-JAR):
	@mkdir -p $(dir $@)
	curl --fail --location --silent --show-error \
	  '$(YAML-PARSER-URL)' -o '$@.tmp'
	mv '$@.tmp' '$@'

ifdef YAML_PARSER_DIR
$(YAML-PARSER-SRC-STAMP):
	rm -rf $(YAML-PARSER-SRC-DIR)
	mkdir -p $(YAML-PARSER-SRC-DIR)/yaml_parser
	cp $(YAML_PARSER_DIR)/src/yaml_parser/*.clj* \
	  $(YAML-PARSER-SRC-DIR)/yaml_parser/
	touch $@
else
$(YAML-PARSER-SRC-STAMP): $(YAML-PARSER-JAR)
	rm -rf $(YAML-PARSER-SRC-DIR)
	mkdir -p $(YAML-PARSER-SRC-DIR)
	unzip -oq $< 'yaml_parser/*.clj' 'yaml_parser/*.cljc' \
	  -d $(YAML-PARSER-SRC-DIR)
	touch $@
endif

$(SOURCE_CACHE)/yaml_parser/%.clj: $(YAML-PARSER-SRC-STAMP)
	@mkdir -p $(dir $@)
	$(RM) $@
	cp $(YAML-PARSER-SRC-DIR)/yaml_parser/$*.cljc $@

$(SOURCE_CACHE)/yaml_parser/core.clj: $(YAML-PARSER-SRC-STAMP)
	@mkdir -p $(dir $@)
	$(RM) $@
	cp $(YAML-PARSER-SRC-DIR)/yaml_parser/core.clj $@

format: $(GO)
	$(GO) fmt ./...

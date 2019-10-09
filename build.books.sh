#!/usr/bin/env bash

set -e

docker run -v `pwd`:/source jagregory/pandoc -o aprenda-go-com-testes.pdf --latex-engine=xelatex --variable urlcolor=blue --toc --toc-depth=1 pdf-cover.md \
    gb-readme.md \
    meta/motivacao.md \
    primeiros-passos-com-go/instalacao-do-go.md \
    primeiros-passos-com-go/hello-world.md \
    primeiros-passos-com-go/inteiros.md \
    primeiros-passos-com-go/iteracao.md \
    primeiros-passos-com-go/arrays-e-slices.md \
    primeiros-passos-com-go/structs-methods-and-interfaces.md \
    primeiros-passos-com-go/pointers-and-errors.md \
    primeiros-passos-com-go/maps.md \
    primeiros-passos-com-go/injecao-de-dependencia.md \
    primeiros-passos-com-go/mocks.md \
    primeiros-passos-com-go/concorrencia.md \
    primeiros-passos-com-go/select.md \
    primeiros-passos-com-go/reflection.md \
    primeiros-passos-com-go/sync.md \
    primeiros-passos-com-go/context.md \
    criando-uma-aplicacao/introducao.md \
    criando-uma-aplicacao/http-server.md \
    criando-uma-aplicacao/json.md \
    criando-uma-aplicacao/io.md \
    criando-uma-aplicacao/command-line.md \
    criando-uma-aplicacao/time.md \
    criando-uma-aplicacao/websockets.md \
    perguntas-e-respostas/os-exec.md \
    perguntas-e-respostas/error-types.md \

docker run -v `pwd`:/source jagregory/pandoc -o aprenda-go-com-testes.epub --latex-engine=xelatex --toc --toc-depth=1 title.txt \
    gb-readme.md \
    meta/motivacao.md \
    primeiros-passos-com-go/instalacao-do-go.md \
    primeiros-passos-com-go/hello-world.md \
    primeiros-passos-com-go/inteiros.md \
    primeiros-passos-com-go/iteracao.md \
    primeiros-passos-com-go/arrays-e-slices.md \
    primeiros-passos-com-go/structs-methods-and-interfaces.md \
    primeiros-passos-com-go/pointers-and-errors.md \
    primeiros-passos-com-go/maps.md \
    primeiros-passos-com-go/injecao-de-dependencia.md \
    primeiros-passos-com-go/mocks.md \
    primeiros-passos-com-go/concorrencia.md \
    primeiros-passos-com-go/select.md \
    primeiros-passos-com-go/reflection.md \
    primeiros-passos-com-go/sync.md \
    primeiros-passos-com-go/context.md \
    criando-uma-aplicacao/introducao.md \
    criando-uma-aplicacao/http-server.md \
    criando-uma-aplicacao/json.md \
    criando-uma-aplicacao/io.md \
    criando-uma-aplicacao/command-line.md \
    criando-uma-aplicacao/time.md \
    criando-uma-aplicacao/websockets.md \
    perguntas-e-respostas/os-exec.md \
    perguntas-e-respostas/error-types.md
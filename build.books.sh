#!/usr/bin/env bash

set -e

docker run -v `pwd`:/source jagregory/pandoc -o aprenda-go-com-testes.pdf --latex-engine=xelatex --variable urlcolor=blue --toc --toc-depth=1 pdf-cover.md \
    gb-readme.md \
    outros/motivacao.md \
    primeiros-passos-com-go/instalacao-do-go.md \
    primeiros-passos-com-go/ola-mundo/ola-mundo.md \
    primeiros-passos-com-go/inteiros/inteiros.md \
    primeiros-passos-com-go/iteracao/iteracao.md \
    primeiros-passos-com-go/arrays-e-slices/arrays-e-slices.md \
    primeiros-passos-com-go/estruturas-metodos-e-interfaces/estruturas-metodos-e-interfaces.md \
    primeiros-passos-com-go/ponteiros-e-erros/ponteiros-e-erros.md \
    primeiros-passos-com-go/maps/maps.md \
    primeiros-passos-com-go/injecao-de-dependencia/injecao-de-dependencia.md \
    primeiros-passos-com-go/mocks/mocks.md \
    primeiros-passos-com-go/concorrencia/concorrencia.md \
    primeiros-passos-com-go/select/select.md \
    primeiros-passos-com-go/reflection/reflection.md \
    primeiros-passos-com-go/sync/sync.md \
    primeiros-passos-com-go/contexto/contexto.md \
    criando-uma-aplicacao/introducao.md \
    criando-uma-aplicacao/servidor-http/servidor-http.md \
    criando-uma-aplicacao/json/json.md \
    criando-uma-aplicacao/io/io.md \
    criando-uma-aplicacao/linha-de-comando/linha-de-comando.md \
    criando-uma-aplicacao/time/time.md \
    criando-uma-aplicacao/websockets/websockets.md \
    duvidas-da-comunidade/os-exec/os-exec.md \
    duvidas-da-comunidade/error-types/error-types.md \
    outros/glossario.md \

docker run -v `pwd`:/source jagregory/pandoc -o aprenda-go-com-testes.epub --latex-engine=xelatex --toc --toc-depth=1 title.txt \
    gb-readme.md \
    outros/motivacao.md \
    primeiros-passos-com-go/instalacao-do-go.md \
    primeiros-passos-com-go/ola-mundo/ola-mundo.md \
    primeiros-passos-com-go/inteiros/inteiros.md \
    primeiros-passos-com-go/iteracao/iteracao.md \
    primeiros-passos-com-go/arrays-e-slices/arrays-e-slices.md \
    primeiros-passos-com-go/estruturas-metodos-e-interfaces/estruturas-metodos-e-interfaces.md \
    primeiros-passos-com-go/ponteiros-e-erros/ponteiros-e-erros.md \
    primeiros-passos-com-go/maps/maps.md \
    primeiros-passos-com-go/injecao-de-dependencia/injecao-de-dependencia.md \
    primeiros-passos-com-go/mocks/mocks.md \
    primeiros-passos-com-go/concorrencia/concorrencia.md \
    primeiros-passos-com-go/select/select.md \
    primeiros-passos-com-go/reflection/reflection.md \
    primeiros-passos-com-go/sync/sync.md \
    primeiros-passos-com-go/contexto/contexto.md \
    criando-uma-aplicacao/introducao.md \
    criando-uma-aplicacao/servidor-http/servidor-http.md \
    criando-uma-aplicacao/json/json.md \
    criando-uma-aplicacao/io/io.md \
    criando-uma-aplicacao/linha-de-comando/linha-de-comando.md \
    criando-uma-aplicacao/time/time.md \
    criando-uma-aplicacao/websockets/websockets.md \
    duvidas-da-comunidade/os-exec/os-exec.md \
    duvidas-da-comunidade/error-types/error-types.md \
    outros/motivacao.md


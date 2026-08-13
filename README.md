# WASA Project

Applicazione di messaggistica sviluppata con un backend Go e un frontend Vue.

## Versioni

- [`v1.0.0`](https://github.com/Emanuele-Bruni/WASA_project2024/tree/v1.0.0) conserva la prima versione del progetto così com'era al termine dello sviluppo originale.
- `main` contiene lo sviluppo della seconda versione. Tutti i cambiamenti successivi al tag `v1.0.0` fanno parte del percorso verso V2.

Per consultare V1 localmente senza modificare i branch:

```sh
git switch --detach v1.0.0
```

Per tornare allo sviluppo corrente:

```sh
git switch main
```

## Avvio

L'intera applicazione può essere costruita e avviata con Docker Compose:

```sh
docker compose up --build
```

Una volta avviati i container:

- il frontend è disponibile su <http://localhost:8080>;
- il backend è disponibile su <http://localhost:3000>.

La configurazione del backend si trova in `demo/config.yml`; la definizione dei servizi e delle porte è in `docker-compose.yml`.

## Sviluppo della V2

Le modifiche vengono sviluppate su branch dedicati e integrate in `main` tramite pull request piccole e revisionabili. Quando la seconda versione sarà stabile, verrà identificata dal tag `v2.0.0`.

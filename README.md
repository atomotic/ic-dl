# ic-dl

scarica i file MAG XML degli oggetti risultato di una ricerca nella **Biblioteca Digitale** di http://www.internetculturale.it

## Installazione

    ~ go get -u -v github.com/atomotic/ic-dl

## Utilizzo

    ~ ic-dl
        -all
                search all (*)
        -query string
                query string
        -type string
                filtery by type (eg. 'periodici')

### Esempio

    ~ ic-dl -query "archiginnasio"
    oai:www.internetculturale.sbn.it/Teca:20:NT0000:BOA0140	oai-www-internetculturale-sbn-it-teca-20-nt0000-boa0140.xml
    oai:www.internetculturale.sbn.it/Teca:20:NT0000:BOA0010	oai-www-internetculturale-sbn-it-teca-20-nt0000-boa0010.xml
    oai:www.internetculturale.sbn.it/Teca:20:NT0000:BOA0030	oai-www-internetculturale-sbn-it-teca-20-nt0000-boa0030.xml
    [...]

L'output dello script riporta l'identificativo OAI dell'oggetto e il corrispondente file xml salvato.  
Gli xml vengono salvati nella directory `./ic-data`. La directory viene creata automaticamente se non esistente.

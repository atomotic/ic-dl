# ic-magxml-download

cerca una stringa nella **Biblioteca Digitale** di http://www.internetculturale.it e salva i file MAG XML dei risultati.

## installazione

    ~ go get -u -v github.com/atomotic/ic-magxml-download

## utilizzo

    ~ ic-magxml-download
    -query string
    	string to search in Internet Culturale

    ~ ic-magxml-download -query "archiginnasio"
    oai:www.internetculturale.sbn.it/Teca:20:NT0000:BOA0140	oai-www-internetculturale-sbn-it-teca-20-nt0000-boa0140.xml
    oai:www.internetculturale.sbn.it/Teca:20:NT0000:BOA0010	oai-www-internetculturale-sbn-it-teca-20-nt0000-boa0010.xml
    oai:www.internetculturale.sbn.it/Teca:20:NT0000:BOA0030	oai-www-internetculturale-sbn-it-teca-20-nt0000-boa0030.xml
    [..more]

gli xml sono salvati nella directory `./ic-data`. La directory viene creata automaticamente se non esistente.

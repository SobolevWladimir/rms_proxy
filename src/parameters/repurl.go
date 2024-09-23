package parameters

import "net/http"

type ReplacedItem struct {
	path             string              // url который подменяем
	replaceByFakeRms bool                // подменить запрос с помощью другой rms
	fakeRms          RMSConnectParameter // само фейковое рмs
	fakeContent      string              //  ответ .. если не берем из  рмс
}



func(rm *ReplacedItem) Handle( r *http.Request) (string,error) {

	return "", nil;
}

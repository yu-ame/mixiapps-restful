package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/yu-ame/mixiapps-restful/pkg/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"text/template"
)

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathRegExp, _ := regexp.Compile("^/")
	tmplPath := pathRegExp.ReplaceAllString(path, "")
	if tmplPath == "" {
		tmplPath = "index"
	}

	//debug的な表示
	log.Printf("path=" + tmplPath)
	log.Printf("method=" + r.Method)
	log.Printf("requestURL=" + r.RequestURI)
	log.Printf("query=" + r.URL.RawQuery)
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(r.Body)
	body := bufbody.String()
	log.Printf("body=" + body)
	authorization := r.Header.Get("Authorization")
	log.Printf("Authorization=" + authorization)

	contentType := ""
	tmplFullPath := ""
	tmplData := make(map[string]interface{})
	tmplData["AppUrl"] = config.GetString("app_url")

	switch tmplPath {
	case "index":
		contentType += "text/html"
		tmplFullPath += "web/tpl/index.tpl"
	case "xml/gadget.xml":
		contentType += "application/xml"
		tmplFullPath += "web/xml/gadget.xml"
	case "makeRequest":
		verifyError := verifyMakeRequest(r)
		contentType += "text/plain"
		if verifyError != nil {
			tmplFullPath += "web/tpl/mixiapps/restful/verify_error.tpl"
		} else {
			tmplFullPath += "web/tpl/mixiapps/restful/verify_success.tpl"
		}
	default:
		contentType += "text/html"
		tmplFullPath += "web/tpl/notfound.tpl"
	}

	w.Header().Set("Content-Type", contentType)
	tpl := template.Must(template.ParseFiles(tmplFullPath))
	err2 := tpl.Execute(w, tmplData)
	if err2 != nil {
		panic(err2)
	}

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+config.GetString("listen_port"), nil)
}

func verifyMakeRequest(r *http.Request) error {
	//sortめんどくさすぎない
	v := r.URL.Query()
	keys := make([]string, 0, len(v)-1)
	for k := range v {
		if k == "oauth_signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		log.Println(k, v[k])
	}

	//basestring作成
    baseString := ""
	baseString += "GET&" + url.QueryEscape(config.GetString("app_url") + r.URL.Path) + "&"
	for pos, k := range keys {
		if pos != 0 {
			baseString += url.QueryEscape("&")
		}
		baseString += url.QueryEscape(k + "=" + url.QueryEscape(v[k][0]))
	}
	log.Println(baseString)

	//signature
	log.Println("oauth_signature=" + v["oauth_signature"][0])

	// read public key
	publicKey, err := readPublicKey("./web/public-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d", publicKey.E)

	verifyErr := verifySignature(baseString, publicKey, v["oauth_signature"][0])
	if verifyErr != nil {
		return fmt.Errorf("Verification failed: %s", verifyErr.Error())
	}
	log.Println("Congratulations! verity!")
	return nil
}

func readPublicKey(path string) (*rsa.PublicKey, error) {
	publicKeyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyData)
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	return rsaPublicKey, nil
}

func verifySignature(message string, publickey *rsa.PublicKey, base64signature string) error {

	h := crypto.SHA1.New()
	h.Write([]byte(message))
	hashed := h.Sum(nil)

	signature, err := base64.StdEncoding.DecodeString(base64signature)
	if err != nil {
		return err
	}

	err = rsa.VerifyPKCS1v15(publickey, crypto.SHA1, hashed, signature)
	if err != nil {
		return err
	}
	return nil
}

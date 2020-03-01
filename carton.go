// -*- go -*-

package carton

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"

	"compress/gzip"
	"encoding/ascii85"
	"io/ioutil"
)

// file stores all cartonFS Files information.
type file struct {
	// content is an ascii85 encoded string of the gzipped original file
	// content.
	content string
	// modTime is the original file last modification time in nanoseconds.
	modTime int64
}

// cartonFS is a map abstraction of the embeded resources. The map keys is the
// file path and the value is a *file struct.
type cartonFS map[string]*file


// Files returns a list of all files embedex in cartonFS
func (b *cartonFS) Files() []string {
	f := make([]string, len(*b))
	i := 0
	for k := range *b {
		f[i] = k
		i++
	}
	return f
}

// isLocalRecent compares the ModTime of local version and File
// file. Returns true only if local file is more recent than the cartonFS.
//
// If local file does not exists, false is returned.
func (b *cartonFS) isLocalRecent(path string, cartonfile *file) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	if fi.ModTime().UnixNano() > cartonfile.modTime {
		return true
	}
	return false
}

// getFileLocal returns the local file content intead of cartonFS embeded file.
func (b *cartonFS) getFileLocal(path string) (ret []byte, err error) {
	fh, err := os.Open(path)
	if err != nil {
		return
	}
	defer fh.Close()

	var buf bytes.Buffer

	_, err = io.Copy(&buf, fh)
	if err != nil {
		return
	}

	ret = buf.Bytes()
	return
}

// getFileLocal returns the cartonFS file content. However is a local file
// exists and is more recent this later's content is returned.
func (b *cartonFS) getFileFromCarton(path string) (ret []byte, err error) {
	f, ok := (*b)[path]
	if !ok {
		return nil, errors.New("File " + path + " not found.")
	}

	if b.isLocalRecent(path, f) {
		return b.getFileLocal(path)
	}

	// Replace back all tilde chars with backtick.
	decoder := ascii85.NewDecoder(
		strings.NewReader(
			strings.ReplaceAll(f.content, "~", "`")))
	gz, err := gzip.NewReader(decoder)
	if err != nil {
		return
	}
	ret, err = ioutil.ReadAll(gz)
	gz.Close()
	if err != nil {
		return
	}

	return
}

// GetFile return path files content. First try the cartonFS then local storage.
func (b *cartonFS) GetFile(path string) (ret []byte, err error) {
	ret, err = b.getFileFromCarton(path)
	if err != nil {
		ret, err = b.getFileLocal(path)
	}
	return
}


// Begin of dynamic content


// carton conatains the carton data.
var carton = &cartonFS{

	`carton.tpl`: &file{
		content: `
+,^C)z!rn3[6#YO:'0V8g@\Q4RakfkOTe-TR0$q+N>9HWJ2a&$J9h@](YnALI>UQ"pqV-oO.7G2fFY:9
-aa#]6k&H,U"g?3X6"fDBYu*>>d30Kj^jf/;e&7Ul7UnRq/Bd@38!gP57r+(n?cpa!$/a$4V)t#"QK=)
D7Z5Plh+'=,'A]$Vl)kZ%;<lo>UN*-ShQ&hQ-"BpY;Y:=@)iG1oha8Wi2u<@]\5kaR:g+c!#d93LbTB<
E%ICh!E+E*O6"@8lm6N~D)Nll^hQ_bO+rCW8Ma4P,fRgB.;d/eoDN2M92WI6LYW8*ESJmSp&B]heoYPl
7TqlEiUu9-n%9dH*YVE_3S7N(ap>:8:VNCe?/d2c<(dZZs>L<SD'Jp:W;CWHQg"Q]OLNcr2Vg!B).B8*
H;>;^Kh$2U=dI/faT<5%bA1m9T1#Botn%pdpU5Ln5kkA*k6EcoWTVt[Bm)@=8:pMkmPC)+N;8_)";F[=
RK\O+<N[)^m>[5,gG$)th)sC-8B6@@_9[W;.=XK,^*IV.c4+sr/9J^spmbe_$/BcQ=3Qib8~@+#!@[2:
M)U[GM6_u&F$Bc6'nnY7XOjMPpGX1SNDW6U!SXJPXP>_-uO'm'Z!&Y]!16mL.=9KT!poVoZ~LJr$<CLT
9S630'*GFI65:G^N#,LHg!/jEVbF^!K,c)\aH0-\"$K.V*-.Ks,JRu4i[$DOcNk]R?6rTg^L7K-ITN*t
~E(Jpj;O%BaGTs4;q4Ci=GUp-M9)*;cPg@caJIlB?KaFEW3$!>)3Ma"&Mg/Hg[+oE~;7'<(7BH*/#d0e
/C5otTWeA6~LYFn*CQXt1f?1nYG0SNT,O4HJ$L+rHd=U0EX55fe9#e.A@A2ON#gJRJhFi]kZ0\ZMIgO9
e.l/6(Qh=[&5X+?T:^EBja^5]l;9l7,a6,eBRSjV[S+l].2R\<8S.ekdA@Mu-~2t$O_P@(%>PHqQe81$
/g7$tM/:Mk)2IOD@VhT.~<d0+l(M:.BR9'85=ulQth\S;~K,i0H)a&m<gb~a=]&e2Tq,@uArqJNk5aH3
PRQhE_!iRTY_6Fdc67nX3BS^!f/OHcYf"l^p[pd=XB#dIg1lR$pdX7csIu(-/bH~rqC%V[!>$P*lYODA
&O??CTJc+ta"8R@O;SVQ^0']8N[fT19T@:&6s%/h>TfLeibO/lU?NgA@=?e8+Vat2.8ZYLR*/bZ;Q$".
\4e+"QDc%m!ddo_G:LEeUi/;0S8olHMBBoE($m4/N-AtmdPko1^<#FT]6&S#dF2i3dhHY!<5S#L$&72<
!3t&1[Jhq-fl_Bg9EmFX)NlNskGoY>_~3T:8R.FOLVY^aG8RtctanZ+qP8V,-pNN(YrJ"nTgZ0P<Rg8"
2D-/]5a8g?UB2~P0lLq3qrWrQ)s8PHY9hJb"!!!
`,
		modTime: 1583192526000000000,
	},

}


// End of dynamic content

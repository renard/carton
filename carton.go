// Package carton is a ambedded resources file
// generated with carton.
//
// See https://github.com/renard/carton
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

// cartonFS is a map abstraction of the embedded resources. The map keys is
// the file path and the value is a *file struct.
type cartonFS map[string]*file

// Files returns a list of all files embedded in cartonFS
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

// getFileLocal returns the local file content intead of cartonFS embedded file.
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
+,^C)z!rn3[6#YO:'0V8g@\Q4Raki*DTe-TR0$q+N>9HWJ-Tr>:9h@](YnALI>UQ"pqV-oO.7G2fFY8R
Uaa#]2bp,_Y"b2>?-TZq@&bu9RRQrkE+GLfh(0e$,bhb8A5cPl.<C)BYXmju[Hf_~~#j[K4XK#t*/L##
E+5K@cFinXs2/g]8=k(<s<BSIE-PJ;5L='*6"j?K5i6gX;bL'L'_~q_rf7a1L>GXEJ1R@-^)+KP0<S^/
&L(N1.<p?b"Y#4JrDdm[^bhfBaB-j<>M%)F_/34>jT4hKEKko56"/Vj!KCFbsq88#0JDE.*N\%4\BT2Y
~I1T(>8d6G#[Lc(7,Je:1W4!:B0NN<gMa"/O6K'JbnR?%]*)-6)%7hoRjWkajiB3Z-;Dk.GnDUGQ/[TN
M(u76uZ3<YnDnfm_TYJh_h&aK?U3>CN-K(\onW=,;]\?mT^;GR*jTpd)L4VnR[ifc3<ZXZD~s@1'dOqN
0j<)<b)[YDUA"gKd",hVH2$*?lGq+[i-,\S;RaMQZPC)+N;8_)";F[IVK\O,gO!Dgn>[5,gG$)th)sC-
8B6@@_9[W90/<e)jNrE7V4+sr/68Nnfmbe_$/BcQ=3Qib8~A0_+AX.UZ)U[GM6bOad$M#$*nloj^OjMP
pGX1RcDpmu5>qb3Xi?iR=G)p$QIfM>2"Yl9p_,pgZe+g&l+iVZQ^5p%#7(!4>B]B+X.-KCjV!)(aM^5O
g:Pj<_W<ce$Fc2"VG\CFj[5!)+=~"/@,mB:gYFQIaPU)Q'kJ\~\(!\EM;E7$B,sa#P"=m+DVcFnr#]5G
3^0~2JZaX!Y,hYb!!bcCFLRrStF;B@gDC^j8CTTaWQC>4KF~C^.rKL:I"Dai+-nK\S]q^?7$56aS1.\f
Y9D/1<qDkAR^D8,!OB[1a)i$l(=S<jge9_M\q;eaQ;PI_Rq^^Z5Jhm8HN/-D(]pB\"c_Umc_BhM=<s2Z
jZ~MIZ\nJP'Pf#0:4[,Nun%pI"#qC6VgkA8d<XsH7dJhIo?et~=D5usVRc,_IIX+"5@)C0d(ZC[SmXe+
cb\lf(XnkbNEr6B\060Yi3Ln-6n/V"<eYLmaQKMi/;b<q~j861a=g##Frl@V:TW6oSa'LS.[:jaB8bB$
$&Xoi+0$\"]nZ(8^1'rr8i2ni0W7W3PC;I2,ND~rn6'2qLeWfBT^+Ua-dVg9.45W;PT9o1XFrqOpdZdr
:V;0DmfV8E,:R]\$i[JI51#O\hJ)79!SG@4?0UX3[,gMFhN7Z35&)J@&*UC25P;h_.6$,n0op"BmD7-+
K[.cA#rYraFEM(qDPkcir<1)Y36&S#d1Z>Q0hct*<&AIcTUP0NLjMjM:]\n_iP4d3qGW@QEZK!YK732^
,1pQC2Ia%U>9[iqD]HjJUqNUjYmO=*%V;"(+m!cjXq:jsS,b*uW1i6>QD(#C,$icir~9\l-5/li[!!!$
!s6j3sq900G!!
`,
		modTime: 1583266572000000000,
	},
}

// End of dynamic content

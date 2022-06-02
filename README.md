# cssrf

<div align=center>
<pre><code>
<b>Extract juicy information using CSS injection
<i>especially csrf token 🥜</i></b>
<br>
<sup>Basically the same thing as <i><a href=https://github.com/d0nutptr/sic>https://github.com/d0nutptr/sic</a></i> but in Golang
<br>I try my best to change the rust code but I lost so many time</sup>
</code></pre>
</div>

## Show me!

All you need is **launch** `cssrf`:
```shell
cssrf [flags] # nothing crazy => cssrf -h to get flags info
```

**Inject** `https://[ATTACKER_URL]/malicious.css` and **wait**:

![demo](https://github.com/ariary/cssrf/blob/main/cssrf.gif)

*This help me solving a [root-me](https://www.root-me.org/fr/Challenges/Web-Client/CSS-Exfiltration) challenge*

<sup>Posting solution is forbidden, thus the csrf token is not entire</sup>

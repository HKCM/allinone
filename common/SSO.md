

## 认证与授权

### Authentication

Authentication(身份验证), the process of verifying that "you are who you say you are", 这个过程是为了证明你是你。通常来说有这么几个方式:

-   用户名密码认证
-   手机和短信验证码认证
-   邮箱和邮件验证码认证
-   人脸识别/指纹识别的生物因素认证
-   OTP(One Time Password) 认证
-   Radius(Remote Authentication Dial-In User Server,远程认证拨号用户服务) 网络认证

### Authorization

Authorization(权限验证), the process of verifying that "you are permitted to do what you are trying to do". 这个过程是为了证明你是否拥有做这件事的权限,比如修改某个表格等等,如果没有权限的话通常会返回 403 错误码。



| 认证                                                         | 授权                                                       |
| ------------------------------------------------------------ | ---------------------------------------------------------- |
| 验证确认身份以授予对系统的访问权限。                         | 授权确定你是否有权访问资源。                               |
| 这是验证用户凭据以获得用户访问权限的过程。                   | 这是验证是否允许访问的过程。                               |
| 它决定用户是否是他声称的用户。                               | 它确定用户可以访问和不访问的内容。                         |
| 身份验证通常需要用户名和密码。                               | 授权所需的身份验证因素可能有所不同,具体取决于安全级别。   |
| 身份验证是授权的第一步,因此始终是第一步。                   | 授权在成功验证后完成。                                     |
| 例如,特定大学的学生在访问大学官方网站的学生链接之前需要进行身份验证。这称为身份验证。 | 例如,授权确定成功验证后学生有权在大学网站上删除哪些信息。 |



## Session

服务器为了保存用户状态而创建的一个特殊的对象.

http协议本身是一种无状态的协议,而这就意味着用户每一次请求都需要提供用户名和密码来进行用户认证,因为根据http协议,服务器并不能知道是哪个用户发出的请求. 所以为了让服务器能识别是哪个用户发出的请求,于是在服务器存储上一份用户登录的信息,这份登录信息就是Session,Session会在服务器响应时传递给浏览器,浏览器将其保存为cookie,下次请求时会带着cookie一并发送给服务器,这样服务器就能识别请求来自哪个用户了,这就是传统的基于session认证。

### Session缺点

-   **Session**: 每个用户经过服务器认证之后,都要在服务端做一次记录,以方便用户下次请求的鉴别,通常而言session都是保存在内存中,而随着认证用户的增多,服务端的开销会明显增大
-   **扩展性**: 用户认证之后,服务端做认证记录,如果认证的记录被保存在内存中的话,这意味着用户下次请求还必须要请求在这台服务器上,这样才能拿到授权的资源,这样在分布式的应用上,相应的限制了负载均衡器的能力。这也意味着限制了应用的扩展能力。
-   **CSRF**: 因为是基于cookie来进行用户识别的, cookie如果被截获,用户就会很容易受到跨站请求伪造的攻击



## Token

基于token的鉴权机制类似于http协议也是无状态的,它不需要在服务端去保留用户的认证信息或者会话信息。这就意味着基于token认证机制的应用不需要去考虑用户在哪一台服务器登录了,这就为应用的扩展提供了便利。

流程:

-   用户使用用户名密码来请求服务器
-   服务器进行验证用户的信息
-   服务器通过验证发送给用户一个token
-   客户端存储token,并在每次请求时附送上这个token值
-   服务端验证token值,并返回数据

这个token必须要在每次请求时传递给服务端,它应该保存在请求头里, 另外,服务端要支持`CORS(跨来源资源共享)`策略,一般我们在服务端这么做就可以了`Access-Control-Allow-Origin: *`。



### ID Token

ID Token本质上是一个 [`JWT Token`](#JWT),包含了该用户身份信息相关的 key/value 键值对,例如:

```json
{
   "iss": "https://server.example.com",
   "sub": "24400320", // subject 的缩写,为用户 ID
   "aud": "s6BhdRkqt3",
   "nonce": "n-0S6_WzA2Mj",
   "exp": 1311281970,
   "iat": 1311280970,
   "auth_time": 1311280969,
   "acr": "urn:mace:incommon:iap:silver"
}  
```

**ID Token** 本质上是一个 `JWT Token` 意味着:

-   用户的身份信息直接被编码进了 `id_token`,你不需要额外请求其他的资源来获取用户信息；
-   `id_token` 可以验证其没有被篡改过

#### ID Token 完整字段含义

| 字段名                | 翻译                                    |
| :-------------------- | :-------------------------------------- |
| sub                   | subject 的缩写,唯一标识,一般为用户 ID |
| name                  | 姓名                                    |
| given_name            | 名字                                    |
| family_name           | 姓氏                                    |
| middle_name           | 中间名                                  |
| nickname              | 昵称                                    |
| preferred_username    | 希望被称呼的名字                        |
| profile               | 基础资料                                |
| picture               | 头像                                    |
| website               | 网站链接                                |
| email                 | 电子邮箱                                |
| email_verified        | 邮箱是否被认证                          |
| gender                | 性别                                    |
| birthdate             | 生日                                    |
| zoneinfo              | 时区                                    |
| locale                | 区域                                    |
| phone_number          | 手机号                                  |
| phone_number_verified | 认证手机号                              |
| address               | 地址                                    |
| formatted             | 详细地址                                |
| street_address        | 街道地址                                |
| locality              | 城市                                    |
| region                | 省                                      |
| postal_code           | 邮编                                    |
| country               | 国家                                    |
| updated_at            | 信息更新时间                            |



### Access Token

Access Token 用于基于 Token 的认证模式,允许应用访问一个资源 API。用户认证授权成功后会签发 Access Token 给应用。Access Token 的格式可以是 [JWT](#JWT)也可以是一个随机字符串. 应用需要**携带 Access Token** 访问资源 API,资源服务 API 会通过拦截器查验 Access Token 中的 `scope` 字段是否包含特定的权限项目,从而决定是否返回资源。

**绝对不要**使用 Access Token 做认证。Access Token 本身**不能标识用户是否已经认证**。

如果你的用户通过社交账号登录,例如微信登录,微信作为身份提供商会颁发自己的 Access Token,你的应用可以利用 Access Token 调用微信相关的 API。这些 Access Token 是由社交账号服务方控制的,格式也是任意的。

Access Token 中只包含了用户 id,在 `sub` 字段。在你开发的应用中,应该将 Access Token **视为一个随机字符串**,不要试图从中解析信息。

Access Token 内容示例:

```json
{
  "jti": "YEeiX17iDgNwHGmAapjSQ",
  "sub": "601ad46d0a3d171f611164ce", // subject 的缩写,为用户 ID
  "iat": 1612415013,
  "exp": 1613624613,
  "scope": "openid profile offline_access",
  "iss": "https://yelexin-test1.authing.cn/oidc",
  "aud": "601ad382d02a2ba94cf996c4" // audience 的缩写,为应用 ID
}
```

注意 Access Token 不包含除 id 之外的任何用户信息。包含 scope 权限项目,用于调用受保护的 API 接口。所以 Access Token 用于**调用接口**,**而不是用作用户认证**。



### Refresh Token

AccessToken 和 IdToken 是 [JSON Web Token](#JWT),**有效时间**通常较短。通常用户在获取资源的时候需要携带 AccessToken,当 AccessToken 过期后,用户需要获取一个新的 AccessToken。

**Refresh Token** 用于获取新的 AccessToken。这样可以缩短 AccessToken 的过期时间保证安全,同时又不会因为频繁过期重新要求用户登录。

用户在初次认证时,Refresh Token 会和 AccessToken、IdToken 一起返回。你的应用必须安全地存储 Refresh Token,它的**重要性**和密码是一样的,因为 Refresh Token 能够一直让用户保持登录。

以下是 Token 端点返回的 Refresh Token:

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InIxTGtiQm8zOTI1UmIyWkZGckt5VTNNVmV4OVQyODE3S3gwdmJpNmlfS2MifQ.eyJqdGkiOiJ4R01uczd5cmNFckxiakNRVW9US1MiLCJzdWIiOiI1YzlmNzVjN2NjZjg3YjA1YTkyMWU5YjAiLCJpc3MiOiJodHRwczovL2F1dGhpbmcuY24iLCJpYXQiOjE1NTQ1Mzc4NjksImV4cCI6MTU1NDU0MTQ2OSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBvZmZsaW5lX2FjY2VzcyBwaG9uZSBlbWFpbCIsImF1ZCI6IjVjYTc2NWUzOTMxOTRkNTg5MWRiMTkyNyJ9.wX05OAgYuXeYM7zCxhrkvTO_taqxrCTG_L2ImDmQjMml6E3GXjYA9EFK0NfWquUI2mdSMAqohX-ndffN0fa5cChdcMJEm3XS9tt6-_zzhoOojK-q9MHF7huZg4O1587xhSofxs-KS7BeYxEHKn_10tAkjEIo9QtYUE7zD7JXwGUsvfMMjOqEVW6KuY3ZOmIq_ncKlB4jvbdrduxy1pbky_kvzHWlE9El_N5qveQXyuvNZVMSIEpw8_y5iSxPxKfrVwGY7hBaF40Oph-d2PO7AzKvxEVMamzLvMGBMaRAP_WttBPAUSqTU5uMXwMafryhGdIcQVsDPcGNgMX6E1jzLA",
  "expires_in": 3600,
  "id_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InIxTGtiQm8zOTI1UmIyWkZGckt5VTNNVmV4OVQyODE3S3gwdmJpNmlfS2MifQ.eyJzdWIiOiI1YzlmNzVjN2NjZjg3YjA1YTkyMWU5YjAiLCJub25jZSI6IjIyMTIxIiwiYXRfaGFzaCI6Ik5kbW9iZVBZOEFFaWQ2T216MzIyOXciLCJzaWQiOiI1ODM2NzllNC1lYWM5LTRjNDEtOGQxMS1jZWFkMmE5OWQzZWIiLCJhdWQiOiI1Y2E3NjVlMzkzMTk0ZDU4OTFkYjE5MjciLCJleHAiOjE1NTQ1NDE0NjksImlhdCI6MTU1NDUzNzg2OSwiaXNzIjoiaHR0cHM6Ly9hdXRoaW5nLmNuIn0.IQi5FRHO756e_eAmdAs3OnFMU7QuP-XtrbwCZC1gJntevYJTltEg1CLkG7eVhdi_g5MJV1c0pNZ_xHmwS0R-E4lAXcc1QveYKptnMroKpBWs5mXwoOiqbrjKEmLMaPgRzCOdLiSdoZuQNw_z-gVhFiMNxI055TyFJdXTNtExt1O3KmwqanPNUi6XyW43bUl29v_kAvKgiOB28f3I0fB4EsiZjxp1uxHQBaDeBMSPaRVWQJcIjAJ9JLgkaDt1j7HZ2a1daWZ4HPzifDuDfi6_Ob1ZL40tWEC7xdxHlCEWJ4pUIsDjvScdQsez9aV_xMwumw3X4tgUIxFOCNVEvr73Fg",
  "refresh_token": "WPsGJbvpBjqXz6IJIr1UHKyrdVF",
  "scope": "openid profile offline_access phone email",
  "token_type": "Bearer"
}
```

应用携带 Refresh Token 向 Token 端点发起请求时,每次都会返回**相同的 Refresh Token** 和**新的 AccessToken、IdToken**,直到 Refresh Token 过期。









## JWT



JWT(Json Web Token)是由三段信息构成的,将这三段信息文本用`.`链接一起就构成了Jwt字符串

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
```



### header

jwt的头部承载两部分信息:

-   声明类型,这里是jwt
-   声明加密的算法 通常直接使用 HMAC SHA256

完整的头部就像下面这样的JSON:

```bash
{
  'typ': 'JWT',
  'alg': 'HS256'
}
```

然后将头部进行base64加密(该加密是可以对称解密的),构成了第一部分.

```shell
echo -n "{'typ':'JWT','alg':'HS256'}" | base64
eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9
```

### playload

载荷就是存放有效信息的地方。这个名字像是特指飞机上承载的货品,这些有效信息包含三个部分

-   标准中注册的声明
-   公共的声明
-   私有的声明

**标准中注册的声明** (建议但不强制使用) :

-   **iss**: jwt签发者
-   **sub**: jwt所面向的用户
-   **aud**: 接收jwt的一方
-   **exp**: jwt的过期时间,这个过期时间必须要大于签发时间
-   **nbf**: 定义在什么时间之前,该jwt都是不可用的.
-   **iat**: jwt的签发时间
-   **jti**: jwt的唯一身份标识,主要用来作为一次性token,从而回避重放攻击。

**公共的声明** :
 公共的声明可以添加任何的信息,一般添加用户的相关信息或其他业务需要的必要信息.但不建议添加敏感信息,因为该部分在客户端可解密.

**私有的声明** :
 私有声明是提供者和消费者所共同定义的声明,一般不建议存放敏感信息,因为base64是对称解密的,意味着该部分信息可以归类为明文信息。

定义一个payload:

```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
```

然后将其进行base64加密,得到Jwt的第二部分。

```shell
echo -n '{"sub":"1234567890","name":"John Doe","admin":true}'|base64
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9

echo "eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9" | base64 -d
```

### signature

jwt的第三部分是一个签证信息,这个签证信息由三部分组成:

-   header (base64后的)
-   payload (base64后的)
-   secret

这个部分需要base64加密后的header和base64加密后的payload使用`.`连接组成的字符串,然后通过header中声明的加密方式进行加盐`secret`组合加密,然后就构成了jwt的第三部分。

```js
// javascript
var encodedString = base64UrlEncode(header) + '.' + base64UrlEncode(payload);

var signature = HMACSHA256(encodedString, 'secret'); // TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
```

将这三部分用`.`连接成一个完整的字符串,构成了最终的jwt:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
```

**注意:secret是保存在服务器端的,jwt的签发生成也是在服务器端的,secret就是用来进行jwt的签发和jwt的验证,所以,它就是你服务端的私钥,在任何场景都不应该流露出去。一旦客户端得知这个secret, 那就意味着客户端是可以自我签发jwt了。**

### 应用

一般是在请求头里加入`Authorization`,并加上`Bearer`标注:

```bash
fetch('api/user/1', {
  headers: {
    'Authorization': 'Bearer ' + token
  }
})
```

服务端会验证token,如果验证通过就会返回相应的资源。整个流程就是这样的

-   用户使用账号(手机/邮箱/用户名)密码请求服务器
-   服务器验证用户账号是否和数据库匹配
-   服务器通过验证后发送给客户端一个 JWT Token
-   **客户端存储 Token,并在每次请求时携带该 Token**
-   **服务端验证 Token 值,并根据 Token 合法性返回对应资源**



![image-20220307160540570](../images/image-20220307160540570.png)

### 优点

-   因为json的通用性,所以JWT是可以进行跨语言支持的,像JAVA,JavaScript,NodeJS,PHP等很多语言都可以使用。
-   因为有了payload部分,所以JWT可以在自身存储一些其他业务逻辑所必要的非敏感信息。
-   便于传输,jwt的构成非常简单,字节占用很小,所以它是非常便于传输的。
-   它不需要在服务端保存会话信息, 所以它易于应用的扩展

### 安全相关

-   不应该在jwt的payload部分存放敏感信息,因为该部分是客户端可解密的部分。
-   保护好secret私钥,该私钥非常重要。
-   如果可以,请使用https协议



## OAuth 2.0

**OAuth 2.0** 是一个授权标准协议。如果你希望将自己应用的数据安全地授权给调用方,建议使用 OAuth 2.0。

根据 OAuth 2.0 协议规范,主要有**四个主体**:

-   **授权服务器(Authorization Server)**,在成功验证资源所有者且获得授权后负责颁发 Access Token给客户端的服务器。
-   **资源所有者(Resource Owner)**,能够许可受保护资源访问权限的实体。当资源所有者是个人时,它作为最终用户被提及。
-   **客户端(Client)**,使用资源所有者的授权代表资源所有者发起对受保护资源的请求的应用程序。
-   **资源服务器(Resource Server)**,托管受保护资源的服务器,能够接收和响应使用Access Token对受保护资源的请求。

常见的 OAuth 2.0 授权流程如下:

```
 +--------+                               +---------------+
 |        |--(A)- Authorization Request ->|   Resource    |
 |        |                               |     Owner     |
 |        |<-(B)-- Authorization Grant ---|               |
 |        |                               +---------------+
 |        |
 |        |                               +---------------+
 |        |--(C)-- Authorization Grant -->| Authorization |
 | Client |                               |     Server    |
 |        |<-(D)----- Access Token -------|               |
 |        |                               +---------------+
 |        |
 |        |                               +---------------+
 |        |--(E)----- Access Token ------>|    Resource   |
 |        |                               |     Server    |
 |        |<-(F)--- Protected Resource ---|               |
 +--------+                               +---------------+
```

- 授权码模式(authorization code)
- 简化模式(implicit)
- 密码模式(resource owner password credentials)
- 客户端模式

图中所示的抽象OAuth 2.0流程描述了四个角色之间的交互,包括以下步骤:

-   (A)客户端向从资源所有者请求授权。授权请求可以直接向资源所有者发起(如图所示),或者更可取的是通过作为中介的授权服务器间接发起。
-   (B)客户端收到授权许可,这是一个代表资源所有者的授权的凭据,使用本规范中定义的四种许可类型之一或 者使用扩展许可类型表示。授权许可类型取决于客户端请求授权所使用的方式以及授权服务器支持的类型。
-   (C)客户端与授权服务器进行身份认证并出示授权许可请求访问令牌。
-   (D)授权服务器验证客户端身份并验证授权许可,若有效则颁发访问令牌。
-   (E)客户端从资源服务器请求受保护资源并出示访问令牌进行身份验证。
-   (F)资源服务器验证访问令牌,若有效则满足该请求。

如果你想了解更多的 OAuth 2.0 内容,可以阅读[协议规范 (opens new window)](https://tools.ietf.org/html/rfc6749)。

OAuth 2.0 以及 OIDC 的核心就是**授权服务器**。授权服务器用于**签发 Access Token**。每个授权服务器都有一个唯一的 **Issuer URI** 和**签名密钥**。

## OpenID Connect

OpenID Connect 是基于 OAuth 2.0 的身份认证协议,增加了 **Id Token**。OIDC 也制定了 OAuth 2.0 中未定义部分的规范,例如 scope,服务发现,用户信息字段等。

-   **OpenID Provider**,指授权服务器,负责签发 Id Token。
-   **终端用户**,Id Token 的信息中会包含终端用户的信息。
-   **调用方**,请求 Id Token 的应用。
-   **Id Token** 由 OpenID Provider 颁发,包含关于终端用户的信息字段。
-   **Claim** 指终端用户信息字段。

OIDC 的授权流程与 OAuth 2.0 一样,主要区别在于 OIDC 授权流程中会额外返回 Id Token。

## SAML 2.0

**SAML 是 Security Assertion Markup Language 的简称,是一种基于XML的开放标准协议,用于在身份提供者(Identity Provider简称IDP)和服务提供商(Service Provider简称SP)之间交换认证和授权数据。**

```sequence
Title: SAML Follow
participant UserAgent as u
participant ServiceProvider as SP
participant IdentityProvider as IdP
u -> SP: Access Resource
SP -> SP: Generate SAML request
SP --> u: Redirect to IdP with SAML request
u -> IdP: Go to Idp with SAML request
IdP --> u: Login page back
u -> IdP: Login
IdP -> IdP: Parses SAML request and Authenticate
IdP -> IdP: Generate SAML response
IdP --> u: Redirect to SP ACS with SAML response
u -> SP: Go to SP ACS with SAML response
SP -> SP: Verify SAML response
SP --> u: Return resource

```

<img src="../images/image-20220304140244178.png" alt="image-20220304140244178" style="zoom:50%;" />


1.  用户试图登录 SP 提供的应用。
2.  SP 生成 SAML Request,通过浏览器重定向,向 IdP 发送 SAML Request。
3.  IdP 解析 SAML Request 并将用户重定向到认证页面。
4.  用户在认证页面完成登录。
5.  IdP 生成 SAML Response,通过对浏览器重定向,向 SP 的 ACS 地址返回 SAML Response,其中包含 SAML Assertion 用于确定用户身份。
6.  SP 对 SAML Response 的内容进行检验。
7.  用户成功登录到 SP 提供的应用。

## MFA

多因素认证(Multi Factor Authentication,简称 MFA)是一种非常简单的安全实践方法,它能够在用户名称和密码之外再额外增加一层保护。启用多因素认证后,用户进行操作时,除了需要提供用户名和密码外(第一次身份验证),还需要进行第二次身份验证,多因素身份认证结合起来将为你的帐号和资源提供更高的安全保护。

### 多因素认证方式

-   第三方身份验证器
    -   Microsoft Authenticator
    -   Google Authenticator
-   短信/邮箱验证码
-   图形锁
-   小程序认证
-   生物识别
    -   指纹
    -   人脸
    -   声纹
    -   虹膜
-   自适应多因素认证
    -   用户属性: 例如用户名、密码、用户身份等用户自身的属性和信息；
    -   位置感知: 位置感知分为虚拟位置(IP 地址)和物理位置(国家、地区等)；
    -   请求来源: 对当前用户的请求来源进行判断,如:硬件设备信息、用户当前所在的系统等；
    -   生物识别: 使用用户的生物信息进行识别,如:指纹信息、人脸识别等；
    -   行为分析: 是否来自常用的登录地点、是否多次输入错误密码、用户之前的操作记录等一系列用户行为。



## IAM

- IAM:Identity and Access Managetment,身份和访问管理,或者简称身份管理。
- EIAM: Enterprise Identity & Access Management,针对企业内部员工、合作伙伴、临时人员等提供统一身份认证和权限管理能力的内部产品
- CIAM: Customer Identity & Access Management,针对互联网用户的顾客身份管理

例如,京东企业内部的工程师、快递员的身份管理均属于 EIAM,而电商平台中的卖家、买家的身份管理属于 CIAM。

## 扫码登录

```sequence
participant 手机端 as Phone
participant PC端 as PC
participant 认证服务器 as IdP
PC -> IdP: 发送生成二维码的请求(携带PC消息)
IdP -> IdP: 生成二维码并与设备绑定
IdP --> PC: 返回二维码
PC -> PC: 展示二维码
PC -> IdP: 定时轮询刷新二维码状态直到成功
Phone -> PC: 扫描二维码获取二维码ID
Phone -> IdP: 通过手机端身份信息(toke)以及二维码获取临时token
IdP --> IdP: 验证toke以及二维码返回临时token,二维码状态变更为待确认
PC -> PC: 轮询刷新二维码状态变更为待确认
Phone -> IdP: 携带临时token,确认登陆
IdP -> IdP: 二维码状态变更为已确认
PC -> IdP: 定时轮询刷新二维码
IdP --> PC: 返回二维码状态以及PC端token
PC -> ResourceServer: 携带Token请求资源
```


https://zhuanlan.zhihu.com/p/312224113

https://zhuanlan.zhihu.com/p/66037342

https://juejin.cn/post/6942875180565790756


JWT说明: https://www.jianshu.com/p/576dbf44b2ae

OAuth2.0: https://colobu.com/2017/04/28/oauth2-rfc6749/
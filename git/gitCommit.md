# Commit

git commit示例

## 说明
```
标题行: 50个字符以内, 描述主要变更内容, 按以下形式填写
<type>(<scope>): <subject>
- type 提交的类型
   feat: 新特性, 
   fix: 修改问题, 
   docs: 文档修改, 仅仅修改了文档, 比如 README, CHANGELOG, CONTRIBUTE等等, 
   style: 代码格式修改, 修改了空格、格式缩进、逗号等等, 不改变代码逻辑, 注意不是css修改,  
   refactor: 代码重构, 没有加新功能或者修复 bug, 
   perf: 优化相关, 比如提升性能、体验, 
   test: 测试用例, 包括单元测试、集成测试等, 
   chore: 其他修改, 比如改变构建流程, 修改依赖管理, 增加工具等, 
   revert: 回滚到上一个版本。
- scope: 影响的的范围
   影响的的范围, 可以为空, 如Git、CentOS等模块, 或者全局All。
- subject
   主题, 提交描述
   
主体内容: 更详细的说明文本, 建议72个字符以内。 需要描述的信息包括:
* 为什么这个变更是必须的? 
  它可能是用来修复一个bug, 增加一个feature, 提升性能、可靠性、稳定性等等
* 它如何解决这个问题? 具体描述解决问题的步骤
* 是否存在副作用、风险? 
尾部: 如果需要的话可以添加一个链接到issue地址或者其它文档, 或者关闭某个issue。

注意, 标题行、主体内容、尾部之间都有一个空行！
```



## 示例

```
feat(login): add login method

add google login support

issue: xxxxx
```
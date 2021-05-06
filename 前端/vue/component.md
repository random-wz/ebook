### 1. 组件

组件需要注册后才能使用，注册分为全局注册和局部注册两种方式。全局注册后，任何vue实力都可以使用，全局注册示例代码如下：

```vue
Vue.component('my-component',{
	// 选项
})
```

my-component就是注册的组件自定义标签名称，推荐使用小写加减号分割的形式命名，注册之后就可以用`<my-component></my-component>`的形式来使用组件了，如下：

```html

```


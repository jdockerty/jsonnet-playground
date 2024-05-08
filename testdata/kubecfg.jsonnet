// Example kubecfg usage

// Import kubecfg's extended native functions from its embedded filesystem that
// includes the libsonnet file bundled.
local kubecfg = import 'internal:///kubecfg.libsonnet';

{
  myVeryNestedObj:: {
    foo: {
      bar: {
        baz: {
          qux: 'some-val',
        },
      },
    },
  },
  hasValue: kubecfg.objectHasPathAll($.myVeryNestedObj, 'foo.bar.baz.qux'),
}

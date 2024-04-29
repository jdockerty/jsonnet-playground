// Test evaluation with something more than a simple 'hello world'
local projects = ['gruglb', 'jsonnet-playground', 'squirrel'];

{
  [if true then 'projects']: std.map(function(x) { [x]: { repository: 'github.com/jdockerty/%s' % [x] } }, projects),
}

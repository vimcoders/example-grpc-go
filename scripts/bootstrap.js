const rootPwd = os.getenv("MYSQL_ROOT_PASSWORD");
shell.connect("root@mysql-1:3306", rootPwd);

// 配置每个实例
dba.configureInstance('root:' + rootPwd + '@mysql-1:3306');
dba.configureInstance('root:' + rootPwd + '@mysql-2:3306');
dba.configureInstance('root:' + rootPwd + '@mysql-3:3306');

// 创建集群
const cluster = dba.createCluster("example");

// 添加实例
cluster.addInstance('root:' + rootPwd + '@mysql-2:3306', {recoveryMethod: 'clone'});
cluster.addInstance('root:' + rootPwd + '@mysql-3:3306', {recoveryMethod: 'clone'});
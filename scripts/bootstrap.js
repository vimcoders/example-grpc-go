const rootPwd = "Root@123456";
const adminUser = "clusteradmin";
const adminPwd = "Cluster@123456";

shell.connect("root@mysql-1:3306", rootPwd);

let res = session.runSql("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME='mysql_innodb_cluster_metadata'");
if (res.fetchOne()) {
    print("Cluster already exists, skip init.\n");
    exit(0);
}

const cfgOpts = {
    clusterAdmin: adminUser,
    clusterAdminPassword: adminPwd,
    interactive: false,
    restart: true
};

dba.configureInstance(
    {user:"root", password:rootPwd, host:"mysql-1", port:3306},
    cfgOpts
);
dba.configureInstance(
    {user:"root", password:rootPwd, host:"mysql-2", port:3306},
    cfgOpts
);
dba.configureInstance(
    {user:"root", password:rootPwd, host:"mysql-3", port:3306},
    cfgOpts
);

const cluster = dba.createCluster("demoCluster", {interactive:false});

cluster.addInstance(adminUser + "@mysql-2:3306", {
    password: adminPwd,
    recoveryMethod:"clone",
    interactive:false
});
cluster.addInstance(adminUser + "@mysql-3:3306", {
    password: adminPwd,
    recoveryMethod:"clone",
    interactive:false
});

print("==== Cluster init finished ====\n");
pipeline {
    agent {
        label 'master'
    }

    environment {
	NEXUS = "https://gizmo.kruemel.home/nexus/"
	REPOSITORY = "repository/webtools/nmapservice/"
    }



    options {
        ansiColor('xterm')
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '5')
    }

    stages {
        stage('clean') {
            steps {
               sh 'make clean'
            }
        }
	
        stage('tests') {
            steps {
                sh 'make tests'
            }
        }

        stage('build') {
            steps {
                sh 'make nmapservice'
            }
        }

        stage('deploy') {
            when {
                branch 'master'
            }

            steps {
                withCredentials([usernameColonPassword(credentialsId: '80a834f5-b4ca-42b1-b5c6-55db88dca0a4', variable: 'CREDENTIALS')]) {
                    sh 'curl -k -u "$CREDENTIALS" --upload-file bin/nmapservice "${NEXUS}${REPOSITORY}"/nmapservice'
                }
            }
        }

        stage('poke rundeck') {
            when {
                branch 'master'
            }

            steps {
                script {
                    step([$class: "RundeckNotifier",
                        includeRundeckLogs: true,
                        jobId: "8c822ea8-ef03-419d-95cd-5a2ca7106071",
                        rundeckInstance: "gizmo",
                        shouldFailTheBuild: true,
                        shouldWaitForRundeckJob: true,
                        tailLog: true])
                } 
            }
        }
    }

    post {
         always {
             mail to:"rafi@guengel.ch",
              subject:"${JOBNAME}#${BRANCH_NAME} -- ${currentBuild.currentResult}",
              body:"Refer to ${currentBuild.absoluteUrl}"
         }
    }
}

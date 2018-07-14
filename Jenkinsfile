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
		    sh 'curl -k -u "$CREDENTIALS" --data bin/nmapservice "${NEXUS}${REPOSITORY}"/nmapservice'
		}
	    }
	}
    }
}

# $1 = clickid /// You pass in your new string when you run the script from the go server eg:   ./apk2.sh clickid

# Go to APK directory // Change this to wherever your original apk is. Make sure you keep a backup too just in case.
cd /root/Android/src/droidware/app/target/outputs/apk/release

# Create new folder name "clickid"
mkdir $1

# Copy apk to new folder
cp app-release.apk $1

# Go to new folder
cd $1

# Unzip the copied apk. You may need to install zip/unzip command line tools
unzip -q app-release.apk

# Just echoing the hex version of "clickid" for debugging purposes
echo -n $1 | od -A n -t x1 | sed 's/ /\\x/g' | tr -d '\n'

# Convert clickid to hex format and store in variable
CLICKIDHEX=$(echo -n $1 | od -A n -t x1 | sed 's/ /\\x/g' | tr -d '\n')

# Find and replace hex string in resources.arsc. The format here is s/OLDHEX/NEWHEX/g. You will need to place your old hex in there manually
# or you could use the same script used to make CLICKIDHEX. Perl uses the format \x{AA} for heach hex byte.
perl -pe 's/\x{77}\x{4C}\x{31}\x{37}\x{4C}\x{4B}\x{45}\x{37}\x{55}\x{33}\x{4A}\x{34}\x{47}\x{4B}\x{4A}\x{49}\x{48}\x{4A}\x{56}\x{55}\x{4A}\x{50}\x{50}\x{38}/'"$CLICKIDHEX"'/g' < resources.arsc > resources2.arsc

# Get rid of old resources and rename new one. Then delete the copied apk so all we have in the folder are the proper files for zipping the new apk.
rm resources.arsc
mv resources2.arsc resources.arsc
rm app-release.apk

# Zip the apk and call it clickid.apk
zip -r -q $1.apk *

# Here we change back to the original apk folder and run the signapk script.
# You need the signapk.sh script and your keystore file in here.
# When it is done, you should end up with something like signed-clickid.apk if all goes right.
cd ../
./signapk.sh $1/$1.apk keystore1234.keystore 26w0FENYtUsaOQn3GG6Kil#4E\$Zvmx server

# All done. Serve the file now with the go server and then perform cleanup.

